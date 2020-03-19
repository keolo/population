# Population Service

[![Maintainability](https://api.codeclimate.com/v1/badges/bd189711b340c16cddce/maintainability)](https://codeclimate.com/github/keolo/population/maintainability)

## Project

### Problem

* Process population data so that we can retrieve population metadata for a given zip code
* [Specification](docs/specification.md)

### Solution

Importer -> Server <-> Client

* Create an Importer service to Extract, Transform, and Load population data
* Create a Server service to retrieve population metadata for a given zip code
* Create a Client service to consume Server API
* [Architecture](docs/architecture.md)

## Services

### Importer

__Average Import Time (on my macbook pro): ~2s__

[Importer](services/importer) is a Go service used to concurrently extract, transform, and load data from two csv
datasources (cbsa_to_msa.csv and zip_to_cbsa.csv) into an embeded key-value
store (BoltDB) for later retrival.

I wrote an originial brute force implementation in ruby which took around 30 minutes in runtime. I then optimized the performance by using Golang with concurrency to take the runtime to around 2 seconds on my MacBook Pro.

The Importer service crunches through 40k+ zip codes and retrieves the
correct population metadata for each.

The import process results in a BoltDB database file size of only about 8MB.

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

### Server

__Average Request Processing Time: ~5ms__

[Server](services/server) is an HTTP service written in Go. It retrieves population growth
metadata for a given zip and responds to the following endpoint:

`GET /zip/{zip}`

It is containerized and deployed to Google Cloud Run. The
population database (created via the Importer service) is baked into the image
when the container is built.

The public URL can be found here: https://server-7y3morjijq-uw.a.run.app/zip/90065

### Client

[Client](services/client) is a service written in Ruby that demonstrates consumption of the server
API.

It can be invoked via:

`ruby main.rb 90065`

[Mixpanel::Client](https://github.com/keolo/mixpanel_client) is an even more
robust example of an API client that I've written in Ruby.

## TODO

After this concept has been tested and aproved, I would prioritize the following:

* Add tests
* Add quality checks
* Add CI/CD
