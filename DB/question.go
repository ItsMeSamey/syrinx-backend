package DB

import (
  "time"

  "ccs.ctf/utils"

  "go.mongodb.org/mongo-driver/bson"
)

type Question struct {
  ID         uint16  `bson:"questionID"`
  Question   string `bson:"question"`
  Answer     string `bson:"answer"`
  Points     int    `bson:"points"`
  Hint       string `bson:"hint"`
  HintPoints int    `bson:"hintpoints"`
  Level      int    `bson:"level"`
  Timestamp  int64  `bson:"-"`
}

func QuestionFromID(ID uint16, maxTries byte) (question Question, err error) {
  question, ok := QUESTIONS.Get(ID)
  if ok {
    if time.Now().Unix() - question.Timestamp <= 15 { return }
  }

  err = tryHard(func () error {
    return utils.WithStack(QuestionDB.get(bson.M{"questionID": ID}, &question))
  }, maxTries)

  if err != nil { return }

  question.Timestamp = time.Now().Unix()
  QUESTIONS.Set(question.ID, question)
  return
}

