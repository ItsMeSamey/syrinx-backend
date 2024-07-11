package DB

import (
)

type Question struct {
  ID         ObjID  `bson:"_id,omitempty"`
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

