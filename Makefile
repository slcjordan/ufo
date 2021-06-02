PWD=$(shell pwd)

.network: 
	docker network create ufo-demo
	@touch $@

.postgres: .network
	docker run \
		--detach \
		--name ufo_postgres \
		--rm \
		--env POSTGRES_PASSWORD \
		--env PGDATA=/var/lib/postgresql/data/pgdata \
		--env POSTGRES_USER \
		--network 'ufo-demo' \
		--volume ${PWD}/data:/var/lib/postgresql/data \
		--publish 5432:5432 \
		postgres:13.2
	@touch $@

start-postgres: .postgres
	@echo "postgres is running"

stop-postgres:
	docker stop ufo_postgres
	@rm .postgres

migrate:
	goose -dir ./migrations postgres "postgresql://$$POSTGRES_USER:$$POSTGRES_PASSWORD@127.0.0.1:5432/$$POSTGRES_USER?sslmode=disable" up

sqgen:
	sqgen-postgres tables --database "postgres://$$POSTGRES_USER:$$POSTGRES_PASSWORD@127.0.0.1:5432/$$POSTGRES_USER?sslmode=disable" --overwrite

.pgbouncer: start-postgres .network
	docker run \
		--name ufo_pgbouncer \
		--rm \
		--network 'ufo-demo' \
		--env POSTGRESQL_PASSWORD=$$POSTGRES_PASSWORD \
		--env POSTGRESQL_USER=$$POSTGRES_USER \
		--env POSTGRESQL_HOST=ufo_postgres \
		--env PGBOUNCER_EXTRA_ARGS='--verbose' \
		--volume ${PWD}/pgbouncer:/bitnami/pgbouncer/conf \
		--publish 6432:6432 \
		bitnami/pgbouncer:latest
	@touch $@

start-pgbouncer: .pgbouncer
	@echo "pgbouncer is running"

stop-pgbouncer:
	docker stop ufo_pgbouncer
	@rm .pgbouncer
