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
	cbsaToMSA := cbsaToMSA()
	zipToCBSA := zipToCBSA()

	// var altCBSA string
	// var metadata []string
	// var metroStats []metroStat
	// var ms metroStat

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

	// Write to database.
	// spew.Dump(metroStats)
}

// processMetroStat is the ETL process for processing population statisctics for a gien zip code.
// * Find CBSA from Zip
// * Check for alternate CBSA
// * Retrieve population estimates
// * Persist in embeded datastore
func processMetroStat(row []string, cbsaToMSA [][]string, db *bolt.DB, wg *sync.WaitGroup) {
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

	// metroStats = append(metroStats, ms)

	updateMetroStat(db, ms)

	// fmt.Println(ms)
}

func zipToCBSA() [][]string {
	in, err := os.Open("db/zip_to_cbsa.csv")
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

func cbsaToMSA() [][]string {
	in, err := os.Open("db/cbsa_to_msa.csv")
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

func find(records [][]string, val string, col int) []string {
	for _, row := range records {
		if row[col] == val {
			return row
		}
	}
	return []string{}
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
	db, err := bolt.Open("db/population.db", 0600, &bolt.Options{Timeout: 1 * time.Second})
	if err != nil {
		return nil, fmt.Errorf("could not open db, %v", err)
	}
	err = db.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists([]byte("MetroStat"))
		if err != nil {
			return fmt.Errorf("could not create MetroStat bucket: %v", err)
		}
		return nil
	})
	if err != nil {
		return nil, fmt.Errorf("could not set up buckets, %v", err)
	}
	fmt.Println("DB Setup Done")
	return db, nil
}

func updateMetroStat(db *bolt.DB, ms metroStat) {
	bucketName := []byte("MetroStat")
	key := []byte(ms.ZIP)

	value, err := json.Marshal(ms)
	if err != nil {
		fmt.Println("could not marshal entry json: %v", err)
	}

	// fmt.Println(value)

	// store some data
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

	// retrieve the data
	// err = db.View(func(tx *bolt.Tx) error {
	// 	bucket := tx.Bucket(bucketName)
	// 	if bucket == nil {
	// 		return fmt.Errorf("bucket not found: %q", bucketName)
	// 	}

	// 	val := bucket.Get(key)
	// 	// fmt.Println("val:", string(val))

	// 	var stat metroStat
	// 	err := json.Unmarshal(val, &stat)
	// 	if err != nil {
	// 		fmt.Println("error:", err)
	// 	}
	// 	// fmt.Println("stats:", stats)
	// 	// fmt.Printf("%+v\n", stat)
	// 	spew.Dump(stat)

	// 	return nil
	// })

	// if err != nil {
	// 	log.Fatal(err)
	// }
}
