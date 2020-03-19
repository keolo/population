package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gorilla/mux"
	bolt "go.etcd.io/bbolt"
)

type metroStat struct {
	ZIP        string `json:"zip"`
	CBSA       string `json:"cbsa"`
	MSAName    string `json:"msaName"`
	PopEst2014 string `json:"popEst2014"`
	PopEst2015 string `json:"popEst2015"`
}

func main() {
	// Set up database connection.
	db, err := setupDB()
	if err != nil {
		fmt.Errorf("error opening db", err)
	}
	defer db.Close()

	// Configure router.
	router := mux.NewRouter()
	router.HandleFunc("/", homeHandler)
	router.HandleFunc("/zip/{zip}", func(w http.ResponseWriter, r *http.Request) {
		showMetroStatHandler(w, r, db)
	})

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	// Start server.
	srv := &http.Server{
		Handler: router,
		Addr:    fmt.Sprintf("0.0.0.0:%s", port),
		// Good practice: enforce timeouts for servers you create!
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}
	log.Print("server started")
	log.Fatal(srv.ListenAndServe())
}

func homeHandler(w http.ResponseWriter, r *http.Request) {
	log.Print("received request for /")
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "200 Ok")
}

func showMetroStatHandler(w http.ResponseWriter, r *http.Request, db *bolt.DB) {
	log.Print("received request for /zip")

	vars := mux.Vars(r)
	zip := vars["zip"]

	// // Retrieve population data for a given zip code.
	val := fetch(zip, db)
	fmt.Fprint(w, string(val))
}

// setupDB opens a database connection.
func setupDB() (*bolt.DB, error) {
	db, err := bolt.Open("population.db", 0600, &bolt.Options{Timeout: 1 * time.Second})
	if err != nil {
		return nil, fmt.Errorf("could not open db, %v", err)
	}
	log.Print("db setup done")
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

	log.Print("looking up key:", string(key))

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

	// Check if we found a value.
	if len(val) == 0 {
		// Create a blank metroStat with zip and default cbsa.
		ms := metroStat{
			ZIP:  zip,
			CBSA: "99999",
		}
		j, err := json.Marshal(ms)
		if err != nil {
			fmt.Println("could not marshal entry json: %v", err)
		}
		val = []byte(j)
	} else {
		val = val
	}

	return val
}
