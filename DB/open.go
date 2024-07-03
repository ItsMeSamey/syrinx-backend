package DB

import (
	"time"

	"github.com/boltdb/bolt"
)

func tryAddBucket(name string, db *bolt.DB) error {
	return db.Update(func(tx *bolt.Tx) error {
			_, err := tx.CreateBucketIfNotExists([]byte(name))
			return err
	})
}

func OpenDb(name string) (*bolt.DB, error) {
	db, err := bolt.Open(name, 0600, &bolt.Options{Timeout: 1 * time.Second})
	if err != nil || db == nil {
		return nil, err
	}

	if err := tryAddBucket(userBucket, db); err != nil {
		return nil, err
	}
	if err := tryAddBucket(teamBucket, db); err != nil {
		return nil, err
	}

	return db, nil
}

