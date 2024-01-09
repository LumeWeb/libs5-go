package utils

import bolt "go.etcd.io/bbolt"

func CreateBucket(name string, db *bolt.DB) error {
	err :=
		db.Update(func(tx *bolt.Tx) error {
			_, err := tx.CreateBucketIfNotExists([]byte(name))
			if err != nil {
				return err
			}

			return nil
		})

	return err
}
