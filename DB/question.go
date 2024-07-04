package DB

import (
  "encoding/json"
)

type Question struct {
  QuestionID string `json:"questionID"`
  Question   string `json:"question"`
  Points   int  `json:"points"`
  Answer   string `json:"answer"`
  Hint   string `json:"hint"`
}

const (
  questionBucket = "questions"
)

func (question *Question) Create() error {
  data, err := json.Marshal(question)
  if err != nil {
    return err
  }

  return QuestionDB.addToBucket(questionBucket, []byte(question.QuestionID), data)
}

func (question *Question) Get(id string) error {
  data, err := QuestionDB.getFromBucket(questionBucket, []byte(id))
  if err != nil {
    return err
  }

  return json.Unmarshal(data, question)
}

// func deleteQuestion(id string) error {
//   return db.Update(func(tx *bolt.Tx) error {
//     b := tx.Bucket([]byte("Questions"))
//     return b.Delete([]byte(id))
//   })
// }

