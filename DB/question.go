package DB

import (
  "errors"

  "go.mongodb.org/mongo-driver/bson"
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

func questionFromID(ID int16) (*Question, error) {
  var question Question
  err := QuestionDB.get(bson.M{"questionID": ID}, &question)
  if err != nil {
    return nil, errors.New("QuestionFromID: DB.get error\n" + err.Error())
  }
  return &question, nil
}

func GetQuestionFromIDTryHard(ID int16, maxTries byte) (*Question, error) {
  var tries byte = 0

  getQuestion:
  question, err := questionFromID(ID)
  if err != nil {
    if tries > maxTries {
      return nil, errors.New("GetQuestionFromIDTryHard: Error in DB.exists, Max Tries reached\n" + err.Error())
    }
    tries += 1;
    goto getQuestion
  }

  return question, nil
}

