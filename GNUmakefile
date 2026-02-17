default: build

build:
	go build -v ./...

install: build
	go install -v ./...

test:
	go test -v -count=1 -parallel=4 ./...

testacc:
	TF_ACC=1 go test -v -count=1 -parallel=4 -timeout 120m ./...

generate:
	go generate ./...
	go run github.com/hashicorp/terraform-plugin-docs/cmd/tfplugindocs generate

lint:
	golangci-lint run ./...

.PHONY: default build install test testacc generate lint
