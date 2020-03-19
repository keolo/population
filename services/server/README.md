# Server

## Stack

* BoltDB used as an embeded key-value store

## Development

* `go mod tidy`

## Deployment

Build container image.

`gcloud builds submit --tag gcr.io/population-271520/server`

Deploy image to Cloud Run.

`gcloud run deploy --image gcr.io/population-271520/server --platform managed`
