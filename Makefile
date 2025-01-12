BIN_DIR = ./bin
MAIN = ./cmd/user_groups_api/main.go

ALL_CONTAINERS = $(DATABASE_CONTAINER) $(PGADMIN_CONTAINER) $(KAFKA_ZOO)
CONTAINER_APP = go_app
DATABASE_CONTAINER = host_name
PGADMIN_CONTAINER = pgadmin
KAFKA_ZOO = kafka zookeeper

POSTGRES_USER = postgres
DATABASE_NAME = postgres

MIGRATE_FILE_INIT = ./migrations/001_init.sql
MIGRATE_FILE_2 = ./migrations/002_init.sql
MIGRATE_DOWN = ./migrations/down.sql

ALL_SERVICES = $(DB_SERVICE) $(PGADMIN_SERVICE) $(KAFKA_ZOO)
COMPOSE_FILE = docker-compose.yml
APP_SERVICE = app
DB_SERVICE = db
PGADMIN_SERVICE = pgadmin

COMPOSE_FLAGS = -f $(COMPOSE_FILE)



.PHONY: help clean_none

help:
	@echo
	@echo "	Swagger документация: 	http://localhost:8080/swagger/index.html"
	@echo "	pgadmin: 		http://localhost:5050"
	@echo
	@echo "	make pull"
	@echo "	make build"
	@echo "	make up"
	@echo
	@echo "  info           	- Информация о контейнерах и образах"
	@echo "  pull           	- Установить образы postgres, pgadmin4, kafka"
	@echo "  build          	- Собрать образы"
	@echo "  up             	- Запустить все сервисы"
	@echo "  init_db        	- Инициализация БД"
	@echo "  rebuild        	- Пересобрать все сервисы"
	@echo "  down           	- Остановить и удалить все сервисы"
	@echo "  rebuild_app    	- Пересобрать приложение"
	@echo "  restart_app    	- Перезапустить приложение"
	@echo "  restart_db     	- Перезапустить базу данных"
	@echo "  clean      		- Удалить все контейнеры и образы"
	@echo "  logs_app       	- Логи приложения"
	@echo "  drop_db        	- Откатить БД"
	@echo "  start			- Запустить все контейнеры"
	@echo "  stop			- Остановить все контейнеры"
	@echo "  delete_app     	- Удалить приложение"
	@echo "  test           	- Запустить тесты"



info:
	@echo ""
	docker ps -a
	@echo ""
	docker images

pull:
	docker pull postgres:latest
	docker pull dpage/pgadmin4:latest
	docker pull confluentinc/cp-zookeeper:latest
	docker pull confluentinc/cp-kafka:latest

build:
	@docker-compose build

up:
	@docker-compose $(COMPOSE_FLAGS) up -d

down:
	@docker-compose $(COMPOSE_FLAGS) down

rebuild_app:
	docker rm -f go_app
	docker-compose up --build -d app
	@make clean_none

restart_app:
	@docker-compose $(COMPOSE_FLAGS) restart $(APP_SERVICE)

logs_app:
	@docker-compose $(COMPOSE_FLAGS) logs -f $(APP_SERVICE)

clean: clean_containers clean_images clean_none


load_db: init_db migrate2

init_db:
	@docker exec -i $(DATABASE_CONTAINER) psql -U $(POSTGRES_USER) -d $(DATABASE_NAME) < $(MIGRATE_FILE_INIT)

migrate2:
	@docker exec -i $(DATABASE_CONTAINER) psql -U $(POSTGRES_USER) -d $(DATABASE_NAME) < $(MIGRATE_FILE_2)


start:
	docker start $(ALL_CONTAINERS) $(CONTAINER_APP)

stop:
	docker stop $(ALL_CONTAINERS) $(CONTAINER_APP)


rebuild:
	docker-compose $(COMPOSE_FLAGS) up --build -d
	@make clean_none

	

delete_app:
	docker stop go_app
	docker rm go_app
	docker rmi app:latest

test: test_models test_repo test_services


drop_db:
	@docker exec -i $(DATABASE_CONTAINER) psql -U $(POSTGRES_USER) -d $(DATABASE_NAME) < $(MIGRATE_DOWN)


sql_terminal:
	docker exec -it $(DATABASE_CONTAINER) psql -U $(POSTGRES_USER) -d $(DATABASE_NAME)







clean_none:
	@docker image prune -f > /dev/null 2>&1


run:
	go run $(MAIN)


rebuild_without_db:
	docker-compose up -d zookeeper kafka pgadmin


all_logs:
	docker-compose $(COMPOSE_FLAGS) logs -f


swag:
	swag init -g ./cmd/user_groups_api/main.go


test_models:
	@go test ./internal/models
	
test_repo:
	@go test ./internal/repository/

test_services:
	@go test ./internal/services/


clean_containers:
	docker rm -f $$(docker ps -aq)

clean_images:
	docker rmi -f $$(docker images -q)

clean_volume:
	docker volume prune

