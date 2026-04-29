.PHONY: protoc docker-build-run test k8s-deploy k8s-down k8s-tunnel-pmt k8s-tunnel-db

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

k8s-deploy:
	eval $(minikube docker-env)
	docker build -t pmt-api:local .
	kubectl create secret generic db-secret \
	  --from-literal=POSTGRES_USER=postgres \
	  --from-literal=POSTGRES_PASSWORD=postgres \
	  --from-literal=POSTGRES_DB=pmt_db
	kubectl apply -f k8s/


k8s-down:
	kubectl delete -f k8s/
	kubectl delete secret db-secret
	docker rmi pmt-api:local

k8s-tunnel-pmt:
	kubectl port-forward service/pmt 8080:8080 9090:9090

k8s-tunnel-db:
	kubectl port-forward service/db 5432:5432
