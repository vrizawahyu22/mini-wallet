POSTGRESQL_DATABASE=postgres
POSTGRESQL_URL='postgresql://postgres:postgres_password@localhost:5656/mini_wallet_db?sslmode=disable'

migrateInit:
	migrate create -ext sql -dir db/migration -seq $(name)

migrateUp:
	migrate -path db/migration -database ${POSTGRESQL_URL} -verbose up

migrateUp1:
	migrate -path db/migration -database ${POSTGRESQL_URL} -verbose up 1

migrateDown:
	migrate -path db/migration -database ${POSTGRESQL_URL} -verbose down

migrateDown1:
	migrate -path db/migration -database ${POSTGRESQL_URL} -verbose down 1

createDB:
	createdb -U ${POSTGRESQL_DATABASE} mini_wallet_db

sqlc_generate:
	sqlc generate

start:
	docker compose up --build