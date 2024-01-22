#!/bin/bash

# Allow overriding the golangci-lint version or use 'latest'
version="${GOLANGCI_VERSION:-latest}"

# Set the bin path for the current Go environment
gopath_bin="$(go env GOPATH)/bin"

# Install golangci-lint if it's not already installed locally
if [ -z "$(command -v $gopath_bin/golangci-lint)" ]; then
  go install github.com/golangci/golangci-lint/cmd/golangci-lint@${version}
fi

# Run golangci-lint with github-actions format
"$gopath_bin/golangci-lint" run --config .golangci.yaml ./... $@
