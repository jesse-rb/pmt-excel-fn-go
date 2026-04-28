.PHONY: docker-build-run test

docker-build-run:
	docker-compose up --build

docker-down:
	docker-compose down --remove-orphans

test:
	go test $(if $(PKG),$(PKG),./...) \
	$(if $(RUN),-run $(RUN),) \
	-v -json -cover | tparse -all
