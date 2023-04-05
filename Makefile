postgres:
	docker run --name postgres12 -p 5432:5432 -e POSTGRES_USER=root -e POSTGRES_PASSWORD=secret123 -d postgres:12-alpine
createdb:
	docker exec -it postgres12 createdb --username=root --owner=root simple_bank
dropdb:
	docker exec -it postgres12 dropdb simple_bank
create_migration:
	migrate create -ext sql -dir db/migration -seq $(name)
migrate_up:
	migrate -path ./db/migration -database "postgresql://root:secret123@localhost:5432/simple_bank?sslmode=disable" -verbose up ${version}
migrate_down:
	migrate -path ./db/migration -database "postgresql://root:secret123@localhost:5432/simple_bank?sslmode=disable" -verbose down ${version}
sqlc:
	sqlc generate
create_mock_db:
	mockgen -package mockdb -destination db/mock/store.go github.com/abiyogaaron/simplebank-service/db/sqlc Store
test:
	go test -v -cover ./...
server_run:
	go run main.go

.PHONY: postgres createdb dropdb migrate_up migrate_down sqlc server_run create_mock_db create_migration test