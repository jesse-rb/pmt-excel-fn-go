.PHONY: build-run test

build-run:
	docker-compose up --build

test:
	go test $(PKG) -run $(RUN) -v
