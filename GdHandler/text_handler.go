package GdHandler

import (
  "log"
  "errors"
  "encoding/json"
  
  "ccs.ctf/DB"

  "github.com/gorilla/websocket"
)

/// This will probably handle questioning/answering
func (lobby *Lobby) handleTextMessage(message []byte, conn *websocket.Conn) error {
  log.Println("Got String: ", string(message))
  _question := DB.Question{}
  if err := json.Unmarshal(message, &_question); err != nil {
    return err
  }

  if lobby.Team.IsSolved(_question.ID) == true {
    return errors.New("Solved")
  }

  question, err := DB.GetQuestionFromIDTryHard(_question.ID, MAX_TRIES)
  if err != nil {
    return errors.New(("getQuestion: could not get question from ID\n ") + err.Error())
  }

  if question.Level >= lobby.Team.Level {
    return errors.New("getQuestion: Level mismatch")
  }

  var retval []byte
  if _question.Hint == "true" {
    retval, err = getHint(question, lobby.Team, MAX_TRIES)
  } else if _question.Answer != "" {
    retval, err = checkAnswer(question, lobby.Team, _question.Answer, MAX_TRIES)
  } else {
    retval, err = getQuestion(question)
  }

  if err != nil {
    return err
  }
  return conn.WriteMessage(websocket.TextMessage, retval)
}

func getHint(question *DB.Question, team *DB.Team, maxTries byte) ([]byte, error) {
  hint, err := team.GetHint(question, maxTries)
  if err != nil {
    return nil, errors.New(("getHint: get hint\n ") + err.Error())
  }
  
  retval, err := json.Marshal(struct{Hint string}{hint})
  if err != nil {
    return nil, errors.New(("getHint: json Marshal error\n ") + err.Error())
  }

  return retval, nil
}

func checkAnswer(question *DB.Question, team *DB.Team, Answer string, maxTries byte) ([]byte, error) {
  ok, err := team.CheckAnswer(question, Answer, maxTries)
  if err != nil {
    return nil, errors.New(("checkAnswer error\n ") + err.Error())
  }

  if ok {
    return []byte("{\"correct\":true}"), nil
  } else {
    return []byte("{\"correct\":false}"), nil
  }
}

func getQuestion(question *DB.Question) ([]byte, error) {
  retval, err := json.Marshal(struct{
    Question   string
    Level      int
    Points     int
    HintPoints int
  }{
    question.Question,
    question.Level,
    question.Points,
    question.HintPoints,
  })
  if err != nil {
    return nil, errors.New("getQuestion: json Marshal error\n " + err.Error())
  }

  return retval, nil
}

