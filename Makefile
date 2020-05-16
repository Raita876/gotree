VERSION := v0.6.1
PACKAGE_NAME := gotree


.PHONY: build
build:
	go build -o ./bin/$(PACKAGE_NAME) -ldflags "-X main.version=$(VERSION) -X main.name=$(PACKAGE_NAME)" main.go

.PHONY: test
test:
	go test -v

.PHONY: install
install: build
	chmod 755 ./bin/$(PACKAGE_NAME) && mv ./bin/$(PACKAGE_NAME) /usr/local/bin/

.PHONY: run
run:
	go run main.go .

.PHONY: tag
tag:
	git tag $(VERSION)
	git push origin $(VERSION)
