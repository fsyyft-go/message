INTERNAL_CONFIG_PROTO_FILES=$(shell find internal/config -name *.proto)

.PHONY: config
config:
	protoc --proto_path=. \
	       --proto_path=./api/third_party \
 	       --go_out=paths=source_relative:. \
	       $(INTERNAL_CONFIG_PROTO_FILES)