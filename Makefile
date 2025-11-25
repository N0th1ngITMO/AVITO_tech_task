.PHONY: help build run test clean db-up db-down db-logs

APP_NAME=pr-review-app
DB_NAME=pr_review_db

help:
	@echo 'Usage:'
	@echo '  make <target>'
	@echo ''
	@echo 'Targets:'
	@fgrep -h "##" $(MAKEFILE_LIST) | fgrep -v fgrep | sed -e 's/\\$$//' | sed -e 's/##//'

build:
	docker-compose build app

run:
	docker-compose up -d

stop:
	docker-compose down

restart: stop run

logs:
	docker-compose logs -f app

db-logs:
	docker-compose logs -f db

test:
	go test ./... -v

clean:
	docker-compose down -v
	docker system prune -f

db-shell:
	docker-compose exec db psql -U postgres -d $(DB_NAME)

status:
	docker-compose ps


dev: run logs

rebuild: clean build run