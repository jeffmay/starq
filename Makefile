# Helpful variables used by targets
GOBIN ?= $$(go env GOPATH)/bin

GOFILES = $(shell go list -f '{{ range .GoFiles }}{{ $$.Dir }}/{{ . }}{{ "\n" }}{{ end }}' ./...)

export GOPRIVATE = github.com/jeffmay

starq: go.mod go.sum $(GOFILES)
	go build -v .

.PHONY: install
install: starq
	go install .

.PHONY: local
local: starq
	cp starq ${HOME}/bin/starq

.PHONY: pre-commit
pre-commit: test
	./lint.sh

.PHONY: test
test:
	./test.sh

# Remove a lot of the default rules that a golang project doesn't need
# https://www.gnu.org/software/make/manual/html_node/Suffix-Rules.html
.SUFFIXES:
