build:
	docker compose build

run:
	docker compose up -d

run-attached:
	docker compose up

build-and-run: build run

build-and-run-attached: build run-attached

stop:
	docker compose stop

remove:
	docker compose rm -f

clean: stop remove

setup:
	go mod tidy

.PHONY: build run build-and-run stop remove clean setup run-attached build-and-run-attached
