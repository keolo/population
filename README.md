# Population Service

[![Maintainability](https://api.codeclimate.com/v1/badges/bd189711b340c16cddce/maintainability)](https://codeclimate.com/github/keolo/population/maintainability)

## Project

### What

* Demonstrate **an exponentially expensive** task of extracting,
transforming, loading, serving, and consuming data
* [Specification](docs/specification.md)

### Why

* Demonstrate how to maximize efficiency by using appropriate languages,
  infrastructure, and concurrency
* Demonstrate a simple and delightful way to deliver scaleable software
* Highlight business value by demonstrating computational, and human
  process efficiencies which result in on-demand and extremely low operational
  costs

### How

* This project is organized into three distinct services
  * [Importer](services/importer): Extract, transform, and load population data
    from CSV files into an embedded database
  * [Server](services/server): For a given zip code, retrieve and respond with population metadata
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

> _Arrows signify flow of data_

* [Architecture](docs/architecture.md)

### KPIs / Goals

* [x] Import processing time: < 5 minutes (**~2s actual!**)
* [x] API request processing time: < 100ms (**~4ms actual!**)
* [x] Autoscaling: true (**true actual!**)
* [x] Deployment steps: 3 (**2 actual!**)
* [x] Operational cost: < $10/mo (**$0 actual!**)

## Stack

* **Application**
  * Ruby
  * Go
  * BoltDB
* **Delivery**
  * Github
  * CodeClimate
  * Cloud Build
  * Container Registry
* **Infrastructure**
  * Cloud Run

## TODO

After this proof of concept has been vetted and approved, I would prioritize
the following:

* [ ] Add tests
* [ ] Add CI/CD
* [ ] Refactor code
* [x] Clarify documentation
