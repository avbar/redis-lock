build:
	CGO_ENABLED=0 GOOS=linux go build -a -o ./bin/app github.com/avbar/redis-lock/cmd/app

up:
	docker compose up --force-recreate --build

run-all:
	make build
	make up