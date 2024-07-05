package DB

import (
  "time"

  bolt "go.etcd.io/bbolt"
)

type DBInstance struct {
  db *bolt.DB
}

// All the DB declarations
var (
  UserDB DBInstance
  LobbyDB DBInstance
  QuestionDB DBInstance
) 


func (instance *DBInstance) init(fileName string, buckets []string) error {
  db_ptr, err := bolt.Open(fileName, 0600, &bolt.Options{Timeout: 1 * time.Second})
  if err != nil || db_ptr == nil {
    return err
  }
  instance.db = db_ptr

  return instance.db.Update(func(tx *bolt.Tx) error {
    for _, name := range buckets {
      if _, err := tx.CreateBucketIfNotExists([]byte(name)); err != nil {
        return err
      }
    }
    return nil
  })
}

func InitDB() error {
  err := UserDB.init("2024_ctf_users.db", []string{userBucket, teamBucket, sessionBucket})
  if err != nil {
    return err
  }
  err = LobbyDB.init("2024_ctf_lobbies.db", []string{questionBucket})
  if err != nil {
    return err
  }
  err = QuestionDB.init("2024_ctf_questions.db", []string{questionBucket})
  if err != nil {
    return err
  }
  return nil
}


