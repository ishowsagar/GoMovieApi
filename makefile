include .env 

stop_container :
	@echo "stopping all other docker containers"
	@containers="$$(docker ps -q)"; \
	if [ -n "$$containers" ]; then \
		echo "found and stopped containers"; \
		docker stop $$containers; \
	else \
		echo "no containers are running..."; \
	fi


create_container :
	docker run --name ${DB_CONTAINER_NAME} -p 5432:5432 -e POSTGRES_USER=${USER} -e POSTGRES_PASSWORD=${PASSWORD} -d ${POSTGRES_IMAGE}

create_db :
	docker exec -it ${DB_CONTAINER_NAME} createdb --username=${USER} --owner=${USER} ${DB_NAME}

start_container :
	docker start ${DB_CONTAINER_NAME}

create_migrations :
	sqlx migrate add -r init 

migrate :
	sqlx migrate run --database-url "postgres://${USER}:${PASSWORD}@${HOST}:${DB_PORT}/${DB_NAME}?sslmode=disable"

demigrate :
	sqlx migrate revert --database-url "postgres://${USER}:${PASSWORD}@${HOST}:${DB_PORT}/${DB_NAME}?sslmode=disable"


build :
	if [ -f "${BINARY}" ]; then \
		rm ${BINARY}; \
		echo "DELETED ${BINARY}"; \
	fi
	@echo "building binary⚡⚡..."
	go build -o ${BINARY} cmd/server/main.go
	@echo "binary created in root✅✅"

run : 
	@echo "movies api has started🚀🚀..."
	./${BINARY}

stop : 
	@echo "stopping whole backend server❌❌..."
	@-pkill -SIGTERM -f "./${BINARY}"
	@echo "Server stopped successfully🪧."\\


# helpers

fix :
	go mod tidy
add :
	git add .
	@echo "added all the files to version change stash,uncommited yet!"
commit : add
	@echo "Added commit message - ${msg}"
	@ if [ -z "${msg}" ]; then \
	echo "Error : Please provide commit message to commit the changes to the repo."; \
	exit 1; \
	fi
	@echo "commit with message : ${msg}"
	git commit -m "${msg}" 
