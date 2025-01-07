.PHONY: start
start:
	@clear
	go run .

.PHONY: tests
tests:
	@clear
	@echo "running all tests"
	go test -v ./...

