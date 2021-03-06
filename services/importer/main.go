package main

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"sync"
	"time"

	bolt "go.etcd.io/bbolt"
)

const (
	// CBSA represents the CSV column number.
	CBSA = 0
	// MDIV represents the CSV column number.
	MDIV = 1
	// NAME represents the CSV column number.
	NAME = 3
	// LSAD represents the CSV column number.
	LSAD = 4
	// POPESTIMATE2014 represents the CSV column number.
	POPESTIMATE2014 = 11
	// POPESTIMATE2015 represents the CSV column number.
	POPESTIMATE2015 = 12
)

type metroStat struct {
	ZIP        string `json:"zip"`
	CBSA       string `json:"cbsa"`
	MSAName    string `json:"msaName"`
	PopEst2014 string `json:"popEst2014"`
	PopEst2015 string `json:"popEst2015"`
}

func main() {
	db, err := setupDB()
	if err != nil {
		fmt.Println("error opening db:", err)
	}
	defer db.Close()

	// Load CSV files into memory.
	cbsaToMSA := loadCSV("db/cbsa_to_msa.csv")
	zipToCBSA := loadCSV("db/zip_to_cbsa.csv")

	var wg sync.WaitGroup

	// Loop through each zip.
	for _, row := range zipToCBSA {
		// If the CBSA is 99999, the zip code is not part of a CBSA.
		if row[1] == "99999" {
			continue
		}

		// Find each MetroStat concurrently.
		wg.Add(1)
		go processMetroStat(row, cbsaToMSA, db, &wg)
	}

	wg.Wait()

	log.Print("data imported successfully")
}

// processMetroStat is the ETL process for processing population data for a
// given zip code.
func processMetroStat(
	row []string,
	cbsaToMSA [][]string,
	db *bolt.DB,
	wg *sync.WaitGroup,
) {
	defer wg.Done()

	// Populate struct.
	ms := metroStat{
		ZIP:  row[0],
		CBSA: row[1],
	}

	// Check for alternate CBSA.
	altCBSA := findAlternateCBSA(cbsaToMSA, ms.CBSA)

	// Update our CBSA if we've found an alternate.
	if altCBSA != "" {
		ms.CBSA = altCBSA
	}

	// Find metadata for CBSA.
	metadata := findMetadata(cbsaToMSA, ms.CBSA)
	if len(metadata) != 0 {
		// fmt.Println("Found metadata:", metadata)
		ms.MSAName = metadata[NAME]
		ms.PopEst2014 = metadata[POPESTIMATE2014]
		ms.PopEst2015 = metadata[POPESTIMATE2015]
	}

	// Persist to database.
	updateMetroStat(db, ms)
}

func loadCSV(filePath string) [][]string {
	in, err := os.Open(filePath)
	if err != nil {
		log.Fatal(err)
	}

	r := csv.NewReader(in)

	records, err := r.ReadAll()
	if err != nil {
		log.Fatal(err)
	}

	return records
}

func findAlternateCBSA(records [][]string, cbsa string) string {
	for _, row := range records {
		if row[MDIV] == cbsa {
			return row[0]
		}
	}
	return ""
}

func findMetadata(records [][]string, cbsa string) []string {
	for _, row := range records {
		if row[CBSA] == cbsa && row[LSAD] == "Metropolitan Statistical Area" {
			return row
		}
	}
	return []string{}
}

func setupDB() (*bolt.DB, error) {
	db, err := bolt.Open(
		"../server/population.db",
		0600,
		&bolt.Options{Timeout: 1 * time.Second},
	)
	if err != nil {
		return nil, fmt.Errorf("could not open db, %v", err)
	}
	log.Print("db setup done")
	return db, nil
}

func updateMetroStat(db *bolt.DB, ms metroStat) {
	bucketName := []byte("MetroStat")
	key := []byte(ms.ZIP)

	value, err := json.Marshal(ms)
	if err != nil {
		fmt.Println("could not marshal entry json: %v", err)
	}

	// Store the metadata for the given zip/key.
	err = db.Batch(func(tx *bolt.Tx) error {
		bucket, err := tx.CreateBucketIfNotExists(bucketName)
		if err != nil {
			return err
		}

		err = bucket.Put(key, value)
		if err != nil {
			return err
		}
		return nil
	})

	if err != nil {
		log.Fatal(err)
	}
}
