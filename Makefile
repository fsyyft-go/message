API_PROTO_FILES=$(shell find api -not -path "api/third_party/*" -name *.proto)
INTERNAL_CONFIG_PROTO_FILES=$(shell find internal/config -name *.proto)

.PHONY: init
init:
	go install golang.org/x/vuln/cmd/govulncheck@latest
	go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
	go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
	go install github.com/go-kratos/kratos/cmd/kratos/v2@latest
	go install github.com/go-kratos/kratos/cmd/protoc-gen-go-http/v2@latest
	go install github.com/envoyproxy/protoc-gen-validate@latest
	go install github.com/go-kratos/kratos/cmd/protoc-gen-go-errors/v2@latest
	go install github.com/google/gnostic/cmd/protoc-gen-openapi@latest
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	go install github.com/google/wire/cmd/wire@latest

.PHONY: api
api:
	protoc --proto_path=. \
	       --proto_path=./api/third_party \
 	       --go_out=paths=source_relative:. \
 	       --go-http_out=paths=source_relative:. \
 	       --go-grpc_out=paths=source_relative:. \
 	       --openapi_out==paths=source_relative:./api \
	       $(API_PROTO_FILES)

.PHONY: config
config:
	protoc --proto_path=. \
	       --proto_path=./api/third_party \
 	       --go_out=paths=source_relative:. \
	       $(INTERNAL_CONFIG_PROTO_FILES)

.PHONY: generate
generate:
	go mod tidy
	go generate ./...