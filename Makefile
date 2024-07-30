build:
	docker-compose build req_3rd-party_services

run: build
	docker-compose up req_3rd-party_services

test:
	go test -v ./...

migrate-up:
	migrate -path ./migration -database 'postgres://postgres:qwerty@0.0.0.0:5436/postgres?sslmode=disable' up

migrate-down:
	migrate -path ./migration -database 'postgres://postgres:qwerty@0.0.0.0:5436/postgres?sslmode=disable' down