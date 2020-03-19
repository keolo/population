# Server

## Functionality

* BoltDB used as a high performance embeded key-value store
* Database keys are zip codes and values are JSON encoded objects
* Database is embeded in container image during build step

## Benefits

* Scaleable
* Portable
* Maintainable
* Secure
* Low Utilization and Operational Cost

## Local Development

* `cd` to this directory
* `go run main.go`
* Verify server responds to http://0.0.0.0:8080/zip/90065

## Deployment

* Build container image:
 `gcloud builds submit --tag gcr.io/population-271520/server`
* Deploy image to Cloud Run:
`gcloud run deploy --image gcr.io/population-271520/server --platform managed`

## TODO

* Add tests
* Optimize Docker image by using `scratch` instead of `alpine`
* Add CI/CD
