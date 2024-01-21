# Helpful variables used by targets
GOBIN ?= $$(go env GOPATH)/bin

export GOPRIVATE = github.com/jeffmay

.PHONY: test
test:
	if [ -z $$(command -v "${GOBIN}/go-test-coverage") ]; then \
	  go install github.com/vladopajic/go-test-coverage/v2@latest; \
	fi
	go test ./... -coverprofile=./cover.out -covermode=atomic -coverpkg=./...
	${GOBIN}/go-test-coverage --config=./.testcoverage.yaml
