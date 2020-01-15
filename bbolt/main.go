package main

import (
	bolt "go.etcd.io/bbolt"
	"log"
)

func main() {
	db, err := bolt.Open("my.db", 0600, nil)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	if err := db.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists([]byte("test"))
		return err
	}); err != nil {
		log.Println("create bucket is failed", err)
	}
	if err := db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("test"))
		bs := b.Get([]byte("key"))
		log.Println("value is", string(bs))
		return b.Put([]byte("key"), append(bs, []byte("aaaa")...))
	}); err != nil {
		log.Println("update is failed", err)
	}
}
