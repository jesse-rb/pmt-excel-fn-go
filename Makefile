.PHONY: build-run test

build-run:
	docker-compose up --build

test:
	docker-compose run --rm api-test sh -c "\
	go install github.com/mfridman/tparse@v0.18 && \
	go test $(if $(PKG),$(PKG),./...) \
	$(if $(RUN),-run $(RUN),) \
	-v -json | tparse - all"
