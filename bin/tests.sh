#!/bin/bash


# Set the coverage file
COVERAGE_FILE="coverage.out"

# Run go tests with coverage, excluding the mocks directory
go test -coverprofile=$COVERAGE_FILE $(go list ./... | grep -v mocks)

# Check if the coverage file was generated
if [ -f $COVERAGE_FILE ]; then
    # Display coverage details
    go tool cover -func=$COVERAGE_FILE
else
    echo "No coverage file found. Ensure tests are not failing."
fi

