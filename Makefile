build:
	go build -v -o game

psql:
	docker run --name postgres --env POSTGRES_PASSWORD=123456 --env POSTGRES_USER=user --env POSTGRES_DB=tfe --volume postgres-volume:/var/lib/postgresql/data --publish 5432:5432 --detach postgres

migration:
	migrate -path migrations -database "postgres://user:123456@localhost:5432/tfe?sslmode=disable" up

run:
	make psql
	sleep 5
	make build
	make migration
	./game
