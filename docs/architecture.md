# System Architecture

* Assuming that population stats don't change very often, this solution is optimized for a low-write/high-read usecase
* Extract and transform data during import process
* Embeded storage is used for fast access, edge computing, offline access, etc.
* Data updates can be scheduled or triggered on a recurring basis

## Hosting

* Google Cloud Run
* Google Cloud Build

## Services

* Import
  * cbsa(zip)
  * msa(cbsa)
  * metadata(msa)
* Server
  * GET /zip/{zip}
  * processMetroStat(zip)
* Client
  * GET /zip/{zip}

## Models

* MetroStat
  * zip (index)
  * msa_name
  * pop_est_2014
  * pop_est_2015

## Stack

* Frontend
  * Ruby
* Backend
  * Golang
* Storage
  * BoltDB
