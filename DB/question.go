package DB

import (
)

type Question struct {
  ID         ObjID  `bson:"_id,omitempty"`
  Question   string `bson:"question"`
  Points     int    `bson:"points"`
  Answer     string `bson:"answer"`
  Hint       string `bson:"hint"`
}

func QuestionFromID(_id string) (*Question, error) {
  var question Question
  return &question, QuestionDB.get("_id", _id, &question)
}

