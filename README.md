# Population Service

[![Maintainability](https://api.codeclimate.com/v1/badges/bd189711b340c16cddce/maintainability)](https://codeclimate.com/github/keolo/population/maintainability)

## Project

### What

This project demonstrates **an exponentially expensive** task of extracting,
transforming, loading, and consuming data. It consists of three separate
services:

* [Importer](services/importer): Extract, transform, and load population data
  from CSV files to database
* [Server](services/server): Retrieve and respond with population metadata for
  a given zip code
* [Client](services/client): Consume Server API

```text
+------------+        +------------+        +------------+
|            |        |            |        |            |
|            |        |            |        |            |
|  Importer  --------->   Server   <-------->   Client   |
|            |        |            |        |            |
|            |        |            |        |            |
+------------+        +------------+        +------------+
```

Arrows signify flow of data.

### Why

* Demonstrate how to maximize efficiency by using appropriate languages,
  infrastructure, and concurrency
* Demonstrate a simple and delightful development experience
* Highlight business value by demonstrating computational, and human
  process efficiencies which result in on-demand and extremely low operational
  costs

### KPIs / Goals

* Import processing time: < 5 minutes (**~2s actual!**)
* API request processing time: < 100ms (**~4ms actual!**)
* Autoscaling: true (**true actual!**)
* Deployment workflow: 3 steps (**2 steps actual!**)
* Operational cost: < $10/mo (**$0 actual!**)

### Reference

* [Specification](docs/specification.md)
* [Architecture](docs/architecture.md)

## Services

### Importer

#### Average Import Time (on my macbook pro): ~2s

[Importer](services/importer) is a Go service used to concurrently extract,
transform, and load data from two CSV datasources (cbsa_to_msa.csv and
zip_to_cbsa.csv) into an embedded key-value store (BoltDB) for later retrieval.

I wrote an original brute force implementation in ruby which took around 30
minutes in runtime. I then optimized the performance by using Golang with
concurrency to take the runtime to around 2 seconds on my MacBook Pro.

The Importer service crunches through 40k+ zip codes while retrieving the
correct population metadata for each record.

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

#### Average Request Processing Time: ~5ms

[Server](services/server) is an HTTP service written in Go. It retrieves population
growth metadata for a given zip and responds to the following endpoint:

`GET /zip/{zip}`

It is containerized and deployed to Google Cloud Run. The
population database (created via the Importer service) is baked into the image
when the container is built.

The public URL can be found here: https://server-7y3morjijq-uw.a.run.app/zip/90065

### Client

[Client](services/client) is a service written in Ruby that demonstrates
consumption of the server API.

It can be invoked via:

`ruby main.rb 90065`

[Mixpanel::Client](https://github.com/keolo/mixpanel_client) is an even more
robust example of an API client that I've written in Ruby.

## Stack

* Application
  * Ruby
  * Go
  * BoltDB
* Infrastructure
  * Cloud Run
  * Container Registry
* Workflow
  * Github
  * CodeClimate
  * Cloud Build

## TODO

After this proof of concept has been vetted and approved, I would prioritize
the following:

* Add tests
* Add CI/CD
* Refactor code
* Clarify documentation
