package DB

import (
  "errors"
  "log"
  
  bolt "go.etcd.io/bbolt"
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

/// does everything syncronosly
func (instance *DBInstance) forEachInBucket(bucket string, fn func (key, val []byte) error) error {
  return instance.db.View(func (tx *bolt.Tx) error {
    b := tx.Bucket([]byte(bucket))
    if b == nil {
      return errors.New("forEachInBucket: bucket is nil")
    }
    cursor := b.Cursor()
    for k, v := cursor.First(); k != nil; k, v = cursor.Next() {
      if err := fn(k, v); err != nil {
        log.Println(err.Error())
      }
    }
    return nil
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

