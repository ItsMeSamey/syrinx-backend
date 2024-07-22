package GdHandler

import (
  "log"
  "errors"
  "encoding/json"
  
  "ccs.ctf/DB"

  "github.com/gorilla/websocket"
)

/// This will probably handle questioning/answering
func (lobby *Lobby) handleTextMessage(myIndex byte, message []byte, conn *websocket.Conn) error {
  log.Println("Got String: ", string(message))
  _question := DB.Question{}
  if err := json.Unmarshal(message, &_question); err != nil {
    return err
  }

  lobby.PlayerMutex.Lock()
  team, err := lobby.getTeam(myIndex)
  lobby.PlayerMutex.Unlock()
  if err != nil {
    return errors.New(("handleTextMessage: error getting team\n") + err.Error())
  }

  question, err := DB.GetQuestionFromIDTryHard(_question.ID, MAX_TRIES)
  if err != nil {
    return errors.New(("getQuestion: could not get question from ID\n ") + err.Error())
  }

  if question.Level != team.Level {
    return errors.New("getQuestion: Level mismatch")
  }

  var retval []byte
  if _question.Hint == "true" {
    retval, err = getHint(question, team, MAX_TRIES)
  } else if _question.Answer != "" {
    retval, err = checkAnswer(question, team, _question.Answer, MAX_TRIES)
  } else {
    retval, err = getQuestion(question)
  }

  return conn.WriteMessage(websocket.TextMessage, retval)
}

func (lobby *Lobby) getTeam(myIndex byte) (*DB.Team, error) {
  myTeam := lobby.Lobby.Players[myIndex].TeamID
  for i := range lobby.Lobby.Teams {
    if *(lobby.Lobby.Teams[i].TeamID) == *myTeam {
      return &lobby.Lobby.Teams[i], nil
    }
  }
  return nil, errors.New("getTeam: user's team is not present in this lobby")
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
    Hint       string
    HintPoints int
  }{
    question.Question,
    question.Level,
    question.Hint,
    question.HintPoints,
  })
  if err != nil {
    return nil, errors.New("getQuestion: json Marshal error\n " + err.Error())
  }

  return retval, nil
}

