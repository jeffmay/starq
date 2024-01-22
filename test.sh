#!/bin/bash

# Generate coverage report
go test ./... -coverprofile=cover.out -covermode=atomic -coverpkg=./...
go tool cover -html=cover.out -o=cover.html

# Validate target coverage is met
total="$(go tool cover -func=cover.out | tail -1 | awk '{ sub (/%/, "", $3); print $3 }')"
target="$(cat ./codecov.y*ml | yq '.coverage.status.project.default.target' | awk '{ sub (/%/, "", $1); print $1 }')"
if [[ "$total" < "$(cat ./codecov.yaml | yq '.coverage.status.project.default.target')" ]]; then
  echo "ERROR: Coverage is $total, expected at least $target"
  exit 1
fi
