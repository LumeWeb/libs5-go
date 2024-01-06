package utils

import bolt "go.etcd.io/bbolt"

func CreateBucket(name string, db *bolt.DB, cb func(bucket *bolt.Bucket)) error {
	err :=
		db.Update(func(tx *bolt.Tx) error {
			bucket, err := tx.CreateBucketIfNotExists([]byte(name))
			if err != nil {
				return err
			}

			cb(bucket)

			return nil
		})

	return err
}
