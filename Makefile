.PHONY: init install-gofumpt install-air install-golangci-lint install-pymarkdown install-precommit ci run help

all: help

# Continuous integration
ci: init
	go mod tidy -v

# Lint application
lint:
	golangci-lint run ./...

# Run application and database container
run:
	docker compose up -d postgres
	air

# Initialize the repository for development
init: install-gofumpt install-air install-golangci-lint install-precommit install-pymarkdown
ifeq (,$(wildcard .git/hooks/pre-commit))
	pre-commit install
endif

install-gofumpt:
ifeq (, $(shell which gofumpt))
	echo "Installing gofumpt..."
	go install mvdan.cc/gofumpt@latest
endif

install-air:
ifeq (, $(shell which air))
	echo "Installing air..."
	go install github.com/cosmtrek/air@latest
endif

install-precommit:
ifeq (, $(shell which pre-commit))
	echo "Installing pre-commit..."
	python3 -m pip install pre-commit
endif

install-golangci-lint:
ifeq (, $(shell which golangci-lint))
	echo "Installing golangci-lint..."
	$(shell curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(GOPATH)/bin v1.55.2)
endif

install-pymarkdown:
ifeq (, $(shell which pymarkdown))
	echo "Installing pymarkdown..."
	python3 -m pip install pymarkdownlnt
endif

# Show this help
help:
	@awk '/^#/{c=substr($$0,3);next}c&&/^[[:alpha:]][[:alnum:]_-]+:/{print substr($$1,1,index($$1,":")),c}1{c=0}' $(MAKEFILE_LIST) | column -s: -t
