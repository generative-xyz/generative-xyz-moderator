install-go-tools:
	go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.28.0
	go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.2
	go install github.com/envoyproxy/protoc-gen-validate@v0.6.13
	go install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-openapiv2@latest
	go install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-grpc-gateway@latest


install-osx: install-go-tools
	brew install protobuf


update:
	go mod tidy
	go mod vendor

run-server:
	 go run cmd/main.go app

lint:
	golangci-lint run ./...


gen-proto:
	 go run tools/main.go protoc
