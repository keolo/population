package main

import (
	"fmt"
	"log"
	"time"

	bolt "go.etcd.io/bbolt"
)

func main() {
	// Set up database connection.
	db, err := setupDB()
	if err != nil {
		fmt.Errorf("error opening db", err)
	}
	defer db.Close()

	// Retrieve population data for a given zip code.
	val := fetch("90065", db)
	fmt.Println(string(val))
}

// setupDB opens a database connection.
func setupDB() (*bolt.DB, error) {
	db, err := bolt.Open("db/population.db", 0600, &bolt.Options{Timeout: 1 * time.Second})
	if err != nil {
		return nil, fmt.Errorf("could not open db, %v", err)
	}
	// err = db.Update(func(tx *bolt.Tx) error {
	// 	_, err := tx.CreateBucketIfNotExists([]byte("MetroStat"))
	// 	if err != nil {
	// 		return fmt.Errorf("could not create MetroStat bucket: %v", err)
	// 	}
	// 	return nil
	// })
	// if err != nil {
	// 	return nil, fmt.Errorf("could not set up buckets, %v", err)
	// }
	fmt.Println("DB Setup Done")
	return db, nil
}

// fetch returns the value of a matching key in a given bolt database.
// * Retrieve key/zip from db
// * Return value (json string)
// * Return empty json object if not found
func fetch(zip string, db *bolt.DB) []byte {
	bucketName := []byte("MetroStat")
	key := []byte(zip)
	var val []byte

	fmt.Println("Looking up key:", key)

	err := db.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket(bucketName)
		if bucket == nil {
			return fmt.Errorf("bucket not found: %q", bucketName)
		}
		val = bucket.Get(key)
		return nil
	})
	if err != nil {
		log.Fatal(err)
	}

	if len(val) == 0 {
		val = []byte("Empty!")
	} else {
		val = val
	}

	return val
}
