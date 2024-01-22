# Helpful variables used by targets
GOBIN ?= $$(go env GOPATH)/bin

export GOPRIVATE = github.com/jeffmay

.PHONY: pre-commit
pre-commit: test
	./lint.sh

.PHONY: test
test:
	./test.sh
