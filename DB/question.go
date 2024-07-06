package DB

import (
)

type Question struct {
  QuestionID string `bson:"questionID"`
  Question   string `bson:"question"`
  Points     int    `bson:"points"`
  Answer     string `bson:"answer"`
  Hint       string `bson:"hint"`
}

func GetQuestion(_id string) (*Question, error) {
  var question Question
  return &question, QuestionDB.get("_id", _id, &question)
}

