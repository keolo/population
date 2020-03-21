# Server

**Average Request Processing Time: ~5ms**

Server is an HTTP service written in Go. It retrieves population
growth metadata for a given zip and responds to the following endpoint:

`GET /zip/{zip}`

It is containerized and deployed to Google Cloud Run. The
population database (created via the Importer service) is baked into the image
when the container is built.

The public URL can be found here: https://server-7y3morjijq-uw.a.run.app/zip/90065

## Functionality

* BoltDB used as a high performance embedded key-value store
* Database keys are zip codes and values are JSON encoded objects
* Database is embedded in container image during build step

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
