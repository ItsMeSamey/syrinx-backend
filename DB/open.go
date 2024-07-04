package DB

import (
	"time"

	bolt "go.etcd.io/bbolt"
)

var DBInstance *bolt.DB = nil

func tryAddBucket(name string) error {
	return DBInstance.Update(func(tx *bolt.Tx) error {
			_, err := tx.CreateBucketIfNotExists([]byte(name))
			return err
	})
}

func OpenDb(name string) (error) {
	db, err := bolt.Open(name, 0600, &bolt.Options{Timeout: 1 * time.Second})
	if err != nil || db == nil {
		return err
	}
	DBInstance = db

	if err := tryAddBucket(userBucket); err != nil {
		return err
	}
	if err := tryAddBucket(teamBucket); err != nil {
		return err
	}

	return nil
}

