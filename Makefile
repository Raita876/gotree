VERSION := v0.2.0
PACKAGE_NAME := gotree


.PHONY: build
build:
	go build -o ./bin/$(PACKAGE_NAME) -ldflags "-X main.version=$(VERSION) -X main.name=$(PACKAGE_NAME)" main.go

.PHONY: test
test:
	go test -v

.PHONY: install
install: test build
	chmod 755 ./bin/$(PACKAGE_NAME) && mv ./bin/$(PACKAGE_NAME) /usr/local/bin/
