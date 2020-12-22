TAG?=$(shell cat version)
IMAGE?=gabrielcarpr/monzo-ynab
TAGGED_IMAGE?=$(IMAGE):$(TAG)

.PHONY: all test build image clean

all: build

build: clean test src/build/monzo-ynab image

test:
	cd src; go test ./...

clean:
	rm -rf src/build

image:
	docker build -t $(IMAGE):latest -t $(TAGGED_IMAGE) .
	docker push $(IMAGE):latest
	docker push $(TAGGED_IMAGE)

src/build/monzo-ynab:
	cd src; CGO_ENABLED=0 go build -o build/monzo-ynab .
