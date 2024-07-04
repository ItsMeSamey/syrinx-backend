package DB

import (
  "time"
  "errors"

  bolt "go.etcd.io/bbolt"
)

type DBInstance struct {
  db *bolt.DB
}

// All the DB declarations
var (
  UserDB DBInstance
  QuestionDB DBInstance
) 


func addToBucketInternal(tx *bolt.Tx, bucket string, key []byte, val []byte) error {
  b := tx.Bucket([]byte(bucket))
  if b == nil {
    return errors.New("addToBucketInternal: bucket is nil")
  }
  if err := b.Put(key, val); err != nil {
    return err
  }
  return nil
}
func (instance *DBInstance) addToBucket(bucket string, key []byte, val []byte) error {
  return instance.db.Update(func(tx *bolt.Tx) error {
    return addToBucketInternal(tx, bucket, key, val)
  })
}

func getFromBucketInternal(tx *bolt.Tx, bucket string, key []byte) ([]byte, error) {
  b := tx.Bucket([]byte(bucket))
  if b == nil {
    return nil, errors.New("getFromBucketInternal: bucket is nil")
  }
  val := b.Get(key)
  if val == nil {
    return nil, errors.New("getFromBucketInternal: key not found")
  }
  return val, nil
}
func (instance *DBInstance) getFromBucket(bucket string, key []byte) ([]byte, error) {
  var val []byte
  err := instance.db.View(func(tx *bolt.Tx) error {
    var err error = nil
    val, err = getFromBucketInternal(tx, bucket, key)
    return err
  })
  return val, err
}

func deleteInBucketInternal(tx *bolt.Tx, bucket string, key []byte) error {
  b := tx.Bucket([]byte(bucket))
  if b == nil {
    return errors.New("deleteInBucketInternal: bucket is nil")
  }
  return b.Delete(key)
}
func (instance *DBInstance) deleteInBucket(bucket string, key []byte) error {
  return instance.db.Update(func(tx *bolt.Tx) error {
    return deleteInBucketInternal(tx, bucket, key)
  })
}

func (instance *DBInstance) DoesExist(bucket string, key []byte) (bool, error) {
  isThere := true
  return isThere, instance.db.View(func(tx *bolt.Tx) error {
    b := tx.Bucket([]byte(userBucket))
    if b == nil {
      return errors.New("DoesExist: bucket is nil")
    }
    val := b.Get(key)
    isThere = val == nil
    return nil
  })
}

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
  err = QuestionDB.init("2024_ctf_questions.db", []string{questionBucket})
  if err != nil {
    return err
  }
  return nil
}


