#
build:
	go build ./cmd/...

deploy_ideas:
	gcloud app deploy ./cmd/bark-ideas/app.yaml

deploy_dogs:
	gcloud app deploy ./cmd/bark-dogs/app.yaml

deploy_default:
	gcloud app deploy ./cmd/default/app.yaml

deploy_dispatch:
	gcloud app deploy ./cmd/default/dispatch.yaml

clean:
	rm -f bark-service
