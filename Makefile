VERSION := v0.1.0
PACKAGE_NAME := gotree


.PHONY: build
build:
	go build -o ./bin/$(PACKAGE_NAME) -ldflags "-X main.version=$(VERSION)" main.go

.PHONY: test
test:
	go test -v
