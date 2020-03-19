# README

## Project

* Specifications
* Architecture

## Services

### Importer

Importer is a Go service is used to extract, transform, and load data from two csv
datasources (cbsa_to_msa.csv and zip_to_cbsa.csv) into an embeded key-value
store (BoltDB) for later retrival.

The importer retrieves and shapes the data into the following pseudo-structure:

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

### Server

Server is an HTTP service written in Go. It retrieves population growth
metadata for a given zip and responds to the following endpoint:

`GET /zip/{zip}`

It is containerized and deployed to Google Cloud Run. The
population database (created via the Importer service) is baked into the image
when the container is built.

### Client

Client is a service written in Ruby that demonstrates consumption of the server
API.

It can be invoked via:

`ruby main.rb 90065`
