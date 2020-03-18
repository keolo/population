package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
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

	// Start server.
	log.Print("server started")
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		showMetroStatHandler(w, r, db)
	})
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", port), nil))
}

func showMetroStatHandler(w http.ResponseWriter, r *http.Request, db *bolt.DB) {
	log.Print("received request")
	zip := r.URL.Path[1:]

	// Retrieve population data for a given zip code.
	val := fetch(zip, db)
	fmt.Fprint(w, string(val))
}

// setupDB opens a database connection.
func setupDB() (*bolt.DB, error) {
	db, err := bolt.Open("db/population.db", 0600, &bolt.Options{Timeout: 1 * time.Second})
	if err != nil {
		return nil, fmt.Errorf("could not open db, %v", err)
	}
	fmt.Println("db setup done")
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

	fmt.Println("looking up key:", string(key))

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
