server:
	go run cmd/server/main.go

docker-up:
	docker-compose -f docker-compose.yaml up


migrate-up:
	 migrate -path migration -database "postgres://postgres:postgres@localhost:5435/lab3?sslmode=disable" up

migrate-down:
	 migrate -path migration -database "postgres://postgres:1234@localhost:5432/lab3?sslmode=disable" down
