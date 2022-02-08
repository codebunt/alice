package db

import (
	"log"
	"sync"

	"github.com/boltdb/bolt"
)

var instance *KapowDB
var initerror error
var once sync.Once

type KapowDB struct {
	db *bolt.DB
}

func GetKapowDBInstance() (*KapowDB, error) {
	//singleton
	once.Do(func() {
		// Open the my.db data file in your current directory.
		// It will be created if it doesn't exist.
		db, err := bolt.Open("./kapow.db", 0600, nil)
		if err != nil {
			log.Fatal(err)
			initerror = err
			return
		}
		instance = &KapowDB{db: db}
	})
	return instance, initerror
}

func (kdb *KapowDB) StartWritableTx() (*bolt.Tx, error) {
	return kdb.db.Begin(true)
}

func (kdb *KapowDB) StartReadTx() (*bolt.Tx, error) {
	return kdb.db.Begin(false)
}

func (kdb *KapowDB) ensureBucket(txn *bolt.Tx, name string) (*bolt.Bucket, error) {
	return txn.CreateBucketIfNotExists([]byte(name))
}

func (kdb *KapowDB) Set(txn *bolt.Tx, bucket string, key string, value []byte) error {
	bkt, err := kdb.ensureBucket(txn, bucket)
	if err != nil {
		return err
	}
	if err != nil {
		return err
	} else if err := bkt.Put([]byte(key), value); err != nil {
		return err
	}
	return nil
}

func (kdb *KapowDB) Get(txn *bolt.Tx, bucket string, key string) ([]byte, error) {
	bkt, err := kdb.ensureBucket(txn, bucket)
	if err != nil {
		return nil, err
	}
	buf := bkt.Get([]byte(key))
	return buf, nil
}

func (kdb *KapowDB) Close() {
	println("closing db")
	kdb.db.Close()
	println("closed db")

}
