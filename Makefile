.PHONY: protoc docker-build-run test

protoc-gateway:
	mkdir -p protogen && \
	cd ./protobuf/ && \
	protoc -I . --grpc-gateway_out ../protogen/ \
	--grpc-gateway_opt paths=source_relative \
	$$(find . -name "*.proto")

protoc: protoc-gateway
	cd ./protobuf/ && protoc -I . \
    --go_out ../protogen/ --go_opt paths=source_relative \
    --go-grpc_out ../protogen/ --go-grpc_opt paths=source_relative \
    $$(find . -name "*.proto")

docker-build-run: protoc
	docker-compose up --build

docker-down:
	docker-compose down --remove-orphans

test: protoc
	go test $(if $(PKG),$(PKG),./pkg/... ./cmd/...) \
	$(if $(RUN),-run $(RUN),) \
	-v -json -cover | tparse -all

