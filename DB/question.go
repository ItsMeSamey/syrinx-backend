package main


import (
	"encoding/json"
    "errors"
    "log"
    "net/http"
    "github.com/boltdb/bolt"
)

type Question struct {
	QuestionID string `json:"questionID"`
	Question   string `json:"question"`
	Points     int    `json:"points"`
	Answer   string `json:"answer"`
	Hint     string `json:"hint"`
}
var db *bolt.DB

func initDB() error {
    var err error
    db, err = bolt.Open("questions.db", 0600, nil)
    if err != nil {
        return err
    }
    return db.Update(func(tx *bolt.Tx) error {
        _, err := tx.CreateBucketIfNotExists([]byte("Questions"))
        return err
    })
}

func addQuestion(question Question) error {
    return db.Update(func(tx *bolt.Tx) error {
        b := tx.Bucket([]byte("Questions"))
        encoded, err := json.Marshal(question)
        if err != nil {
            return err
        }
        return b.Put([]byte(question.QuestionID), encoded)
    })
}

func getQuestion(id string) (*Question, error) {
    var question Question
    err := db.View(func(tx *bolt.Tx) error {
        b := tx.Bucket([]byte("Questions"))
        data := b.Get([]byte(id))
        if data == nil {
            return errors.New("question not found")
        }
        return json.Unmarshal(data, &question)
    })
    if err != nil {
        return nil, err
    }
    return &question, nil
}

func deleteQuestion(id string) error {
    return db.Update(func(tx *bolt.Tx) error {
        b := tx.Bucket([]byte("Questions"))
        return b.Delete([]byte(id))
    })
}
