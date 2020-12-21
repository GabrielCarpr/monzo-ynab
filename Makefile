TAG?=$(shell cat ${version})
IMAGE?=gabrielcarpr/monzo-ynab:$(TAG)

.PHONY: test build release

build: test monzo-ynab

test:
	cd src; go test ./...

monzo-ynab:
	cd src; go build -o build/monzo-ynab .

