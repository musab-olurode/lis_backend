include .env
run:
	go build && ./lis_backend

tidy:
	go mod tidy && go mod vendor

migrate:
	cd database/schema && goose postgres ${DB_URL} up && cd ../../

migrate-down:
	cd database/schema && goose postgres ${DB_URL} down && cd ../../

sqlc:
	sqlc generate
