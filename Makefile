.PHONY: test-coverage

test-coverage:
    find . -name go.mod -execdir go test -coverprofile=coverage.out ./... \;
	go tool cover -html=coverage.out -o coverage.html
