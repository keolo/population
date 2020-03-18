# System Architecture

* Assuming that population stats don't change very often, this solution is optimized for a read heavy usecase
* Extract and transform data during import process
* In-memory storage and/or caching at the API layer could be used for fast access

## Hosting

* Heroku

## Services

* Import
  * cbsa(zip)
  * msa(cbsa)
  * metadata(msa)
* API
  * GET /metro_stats/:zip
  * metro_stats(zip)
* Client
  * GET /metro_stats/:zip

## Models

* MetroStat
  * zip (index)
  * msa_name
  * pop_est_2014
  * pop_est_2015

## Stack

* Frontend
  * Ruby client
* Backend
  * Ruby
  * Rails API
  * Rspec
* Storage
  * Postgres

## TODO

* Create rake task to trigger import of csv files
* Import CSVs to database tables
