package GdHandler

import (
  "errors"
  "encoding/json"
  
  "ccs.ctf/DB"
)

const (
  QUESTION_L3N6 = 101
)


func specialQuestionHint_L3N6(question *DB.Question, team *DB.Team) ([]byte, error) {
  hint, err := team.GetHint(question, 5)
  if err != nil {
    return nil, errors.New(("getHint: get hint\n ") + err.Error())
  }

  retval, err := json.Marshal(struct{
    Hint string
    Lol string
  }{
    hint,
    question.Answer,
  })
  if err != nil {
    return nil, errors.New(("getHint: json Marshal error\n ") + err.Error())
  }
  return retval, nil
}

