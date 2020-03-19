# Importer

## Functionality

* Find CBSA from Zip
* Check for alternate CBSA
* Retrieve population estimates
* Store population metadata for each zip in BoltDB

## Local Development

* `cd` to this directory
* `go run main.go`
* Verify `../server/population.db` gets created

## TODO

* Add tests
* Add this service to automated build process
* Trigger build and deployment when new data/csv files get uploaded
