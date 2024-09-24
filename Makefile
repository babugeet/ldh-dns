hello:
	echo "Hello"
lint:
	golangci-lint run
fmt:
	cd pkg;go fmt

build: lint fmt 
	cd pkg;go build 
