package db

import (
	"errors"
	"go.etcd.io/bbolt"
)

var _ KVStore = (*BboltDBKVStore)(nil)

type BboltDBKVStore struct {
	db         *bbolt.DB
	bucket     *bbolt.Bucket
	bucketName string
	root       bool
	dbPath     string
}

func (b BboltDBKVStore) Open() error {
	if b.root && b.db == nil {
		db, err := bbolt.Open(b.dbPath, 0666, nil)
		if err != nil {
			return err
		}
		b.db = db
	}

	if b.bucket == nil && len(b.bucketName) > 0 {
		err := b.db.Update(func(txn *bbolt.Tx) error {
			bucket, err := txn.CreateBucketIfNotExists([]byte(b.bucketName))
			if err != nil {
				return err
			}
			b.bucket = bucket
			return nil
		})
		if err != nil {
			return err
		}
	}

	return nil
}

func (b BboltDBKVStore) Close() error {
	if b.root && b.db != nil {
		err := b.db.Close()
		if err != nil {
			return err
		}
	}

	return nil
}

func (b BboltDBKVStore) Get(key []byte) ([]byte, error) {
	if b.root {
		return nil, errors.New("Cannot get from root")
	}

	var val []byte
	err := b.db.View(func(txn *bbolt.Tx) error {
		bucket := txn.Bucket([]byte(b.bucketName))
		val = bucket.Get(key)
		return nil
	})
	if err != nil {
		return nil, err
	}

	return val, nil
}

func (b BboltDBKVStore) Put(key []byte, value []byte) error {

	if b.root {
		return errors.New("Cannot put from root")
	}

	err := b.db.Update(func(txn *bbolt.Tx) error {
		bucket := txn.Bucket([]byte(b.bucketName))
		err := bucket.Put(key, value)
		if err != nil {
			return err
		}
		return nil
	})

	return err
}

func (b BboltDBKVStore) Delete(key []byte) error {
	if b.root {
		return errors.New("Cannot delete from root")
	}

	err := b.db.Update(func(txn *bbolt.Tx) error {
		bucket := txn.Bucket([]byte(b.bucketName))
		err := bucket.Delete(key)
		if err != nil {
			return err
		}
		return nil
	})

	return err
}

func (b BboltDBKVStore) Bucket(prefix string) (KVStore, error) {
	return &BboltDBKVStore{
		db:         b.db,
		bucketName: prefix,
		root:       false,
	}, nil
}

func NewBboltDBKVStore(dbPath string) *BboltDBKVStore {
	return &BboltDBKVStore{
		dbPath: dbPath,
		root:   true,
	}
}
