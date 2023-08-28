.PHONY: vendor

build-cli: vendor
	go build -o ./dist/ ./cmd/deprecated_lambda_func_collector

vendor:
	go mod vendor
	go mod tidy

