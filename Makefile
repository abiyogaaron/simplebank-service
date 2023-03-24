postgres:
	docker run --name postgres12 -p 5432:5432 -e POSTGRES_USER=root -e POSTGRES_PASSWORD=asukuada500 -d postgres:12-alpine
createdb:
	docker exec -it postgres12 createdb --username=root --owner=root simple_bank
dropdb:
	docker exec -it postgres12 dropdb simple_bank
migrate_up:
	migrate -path ./db/migration -database "postgresql://root:asukuada500@localhost:5432/simple_bank?sslmode=disable" -verbose up
migrate_down:
	migrate -path ./db/migration -database "postgresql://root:asukuada500@localhost:5432/simple_bank?sslmode=disable" -verbose down
sqlc:
	sqlc generate
test:
	go test -v -cover ./...

.PHONY: postgres createdb dropdb migrate_up migrate_down sqlc