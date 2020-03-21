# Importer

**Average Import Time (on my macbook pro): ~2s**

Importer is a Go service used to concurrently extract,
transform, and load data from two CSV datasources (cbsa_to_msa.csv and
zip_to_cbsa.csv) into an embedded key-value store (BoltDB) for later retrieval.

> Originally, I wrote a brute force implementation in ruby which took around 30
minutes in runtime. I then optimized the performance by using Golang with
concurrency to take the runtime to around 2 seconds on my MacBook Pro.

* The Importer service crunches through 40k+ zip codes while retrieving the
correct population metadata for each record
* The import process results in a BoltDB database file size of only about 8MB

The importer persists data in the following schema:

```json
{
    "90065": {
        "zip": "90065",
        "cbsa": "31080",
        "msaName": "Los Angeles-Long Beach-Anaheim, CA",
        "popEst2014": "13254397",
        "popEst2015": "13340068"
    },
    "10001": {
        "zip": "10001",
        "cbsa": "35620",
        "msaName": "New York-Newark-Jersey City, NY-NJ-PA",
        "popEst2014": "20095119",
        "popEst2015": "20182305"
    },
    ...
}
```

## Functionality

* Find CBSA from Zip
* Check for alternate CBSA
* Retrieve population estimates
* Store population metadata for each zip in BoltDB

## Benefits

* Fast (~2s on development machine)
* Low Utilization and Operational Cost
  * Because this importer is so efficient, we can run it often at a low cost
* Maintainable

## Local Development

* `cd` to this directory
* `go run main.go`
* Verify `../server/population.db` gets created

## TODO

* Add tests
* Add this service to automated build process
* Trigger build and deployment when new data/csv files get uploaded
* Share `MetroStat` type between Importer and Server code
