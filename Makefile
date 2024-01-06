.PHONY: init install-gofumpt install-golangci-lint install-precommit ci run help

all: help

# Initialize the repository for development
init: install-gofumpt install-golangci-lint install-precommit
ifeq (,$(wildcard .git/hooks/pre-commit))
	pre-commit install
endif

install-gofumpt:
ifeq (, $(shell which golangci-lint))
	echo "Installing gofumpt..."
	go install mvdan.cc/gofumpt@latest
endif

install-pre-commit:
ifeq (, $(shell which golangci-lint))
	echo "Installing pre-commit..."
	python3 -m pip install pre-commit
endif

install-golangci-lint:
ifeq (, $(shell which golangci-lint))
	echo "Installing golangci-lint..."
	$(shell curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(GOPATH)/bin v1.55.2)
endif

# Continuous integration
ci: init
	go mod tidy -v

# Run linter
lint:
	golangci-lint run ./...

# Run application
run:
	air

# Show this help
help:
	@awk '/^#/{c=substr($$0,3);next}c&&/^[[:alpha:]][[:alnum:]_-]+:/{print substr($$1,1,index($$1,":")),c}1{c=0}' $(MAKEFILE_LIST) | column -s: -t
