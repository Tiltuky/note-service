build:
	docker-compose build note-service

run:
	docker-compose up note-service

migrate:
	migrate -path ./schema -database 'postgres://postgres:1qw23er4@0.0.0.0:5436/postgres?sslmode=disable' up

swag:
	swag init -g cmd/main.go
