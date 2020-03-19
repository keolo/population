# Importer

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
