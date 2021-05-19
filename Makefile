PWD=$(shell pwd)

start-postgres:
	docker run \
		--detach \
		--name ufo-demo \
		--rm \
		--env POSTGRES_PASSWORD \
		--env PGDATA=/var/lib/postgresql/data/pgdata \
		--env POSTGRES_USER \
		--volume ${PWD}/data:/var/lib/postgresql/data \
		--publish 5432:5432 \
	postgres:13.2

stop-postgres:
	docker stop ufo-demo

migrate:
	goose -dir ./migrations postgres "postgresql://$$POSTGRES_USER:$$POSTGRES_PASSWORD@127.0.0.1:5432/$$POSTGRES_USER?sslmode=disable" up

sqgen:
	sqgen-postgres tables --database "postgres://$$POSTGRES_USER:$$POSTGRES_PASSWORD@127.0.0.1:5432/$$POSTGRES_USER?sslmode=disable" --overwrite
