all:
.PHONY: all

db-start:
	@echo "starting the postgres docker"
	docker run --rm --name docker-postgres -e POSTGRES_PASSWORD=postgres -d -p 5432:5432 -v ${DEVROOT}/docker/volumes/postgres:/var/lib/postgresql/data postgres:12.3

db-stop:
	@echo "stopping the postgres docker"
	docker container stop docker-postgres

db-ssh:
	@echo "ssh login the docker-postgres container"
	docker exec -it docker-postgres psql -U postgres