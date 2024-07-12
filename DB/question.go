package DB

import (
  "math/rand"
  "errors"
)

type Question struct {
  ID         int16  `bson:"questionID"`
  Question   string `bson:"question"`
  Answer     string `bson:"answer"`
  Points     int    `bson:"points"`
  Hint       string `bson:"hint"`
  HintPoints int    `bson:"hintpoints"`
  Level      int    `bson:"level"`
}

func QuestionFromID(_id string) (*Question, error) {
  var question Question
  return &question, QuestionDB.get("_id", _id, &question)
}

func genQuestionID() (int16, error) {
  times:=0
  start: 
  ID:=int16(rand.Intn(32767))
  exists, err := QuestionDB.exists("questionID", ID)
  if exists {
    if times > 1024 {
      return 0, errors.New("genQuestionID: Lucky Error!!")
    }
    times += 1
    goto start
  }
  return ID, err
}

