#
build:
	go build ./cmd/...

deploy:
	gcloud app deploy ./cmd/bark-service/app.yaml

clean:
	rm -f bark-service
