# Test Coverage
Last Tested: 20240705
Current Results: **85.6%**
- Create file with permissons: `bin/tests.sh`. Give correct permissions (chmod 755)
- To run tests run command `bin/tests.sh` in the terminal.

```
#!/bin/bash

# Set the coverage file
COVERAGE_FILE="coverage.out"

# Run go tests with coverage, excluding the mocks directory
go test -coverprofile=$COVERAGE_FILE $(go list ./... | grep -v mocks | grep -v docs)

# Check if the coverage file was generated
if [ -f $COVERAGE_FILE ]; then
    # Display coverage details
    go tool cover -func=$COVERAGE_FILE
else
    echo "No coverage file found. Ensure tests are not failing."
fi
```
