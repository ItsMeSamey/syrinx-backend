package GdHandler

import (
  "encoding/json"
  "errors"
  // "log"
  
  "ccs.ctf/DB"
  
  "github.com/gorilla/websocket"
)

/// This will probably handle questioning/answering
func (lobby *Lobby) handleTextMessage(message []byte, conn *websocket.Conn) error {
  // log.Println("GOT: ", string(message))
  if LEVEL != lobby.Team.Level && !lobby.Team.Exception {
    return errors.New("getQuestion: Global and Team level mismatch")
  }
  // log.Println("Got String: ", string(message))
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

  if question.Level != lobby.Team.Level && !lobby.Team.Exception {
    return errors.New("getQuestion: Question and Team level mismatch")
  }

  var retval []byte
  if _question.Hint == "true" {
    switch (question.ID) {
    case QUESTION_L3N6:
      retval, err = specialQuestionHint_L3N6(question, lobby.Team)
    default: 
      retval, err = lobby.getHint(question, MAX_TRIES)
    }
  } else if _question.Answer != "" {
    switch (question.ID) {
    default: 
      retval, err = lobby.checkAnswer(question, _question.Answer, MAX_TRIES)
    }
  } else {
    switch (question.ID) {
    default: 
      retval, err = getQuestion(question)
    }
  }

  if err != nil {
    return err
  }
  
  // log.Println("Sent: ", string(retval))
  return conn.WriteMessage(websocket.TextMessage, retval)
}

func (lobby *Lobby) getHint(question *DB.Question, maxTries byte) ([]byte, error) {
  team := lobby.Team

  lobby.PlayerMutex.Lock()
  hint := team.GetHint(question, maxTries)
  lobby.PlayerMutex.Unlock()

  if err := team.Sync(maxTries); err != nil {
    return nil, errors.New(("Team.GetHint: sync error\n ") + err.Error())
  }

  retval, err := json.Marshal(struct{Hint string}{hint})
  if err != nil {
    return nil, errors.New(("getHint: json Marshal error\n ") + err.Error())
  }

  return retval, nil
}

func (lobby *Lobby) checkAnswer(question *DB.Question, Answer string, maxTries byte) ([]byte, error) {
  team := lobby.Team

  lobby.PlayerMutex.Lock()
  if team.IsSolved(question.ID) {
    lobby.PlayerMutex.Unlock()
    return nil, errors.New("Solved")
  }
  ok := team.CheckAnswer(question, Answer, maxTries)
  lobby.PlayerMutex.Unlock()

  if err := team.Sync(maxTries); err != nil {
    return nil, errors.New(("Team.CheckAnswer: sync error\n ") + err.Error())
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

