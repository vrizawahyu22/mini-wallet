# Mini Wallet

## Installation
### Using Docker (Recommended)
If you're using docker, make sure docker has been installed on your computer

### Without Docker
If you're not using docker, there are some tools that you have to install:
  - PostgreSQL
  - Go latest version (I'm using go1.21.3)

## How to Run
### Using Docker (Recommended)
1. Install all dependencies:
```
go mod tidy
```
2. Run App
```
make run
```
if error, please run using this command
```
docker compose up --build
```
### Without Docker
1. Create database using this command
```
make createDB
```
Or you can create the database manually with your DB Management App like PgAdmin

2. Change config in app.env with your database configuration
```
DB_CONN_STRING="dbname=mini_wallet_db user=postgres password=postgres_password host=localhost port=5656 sslmode=disable"
MIGRATION_URL=file://db/migration
```
3. Install all dependencies:
```
go mod tidy
```
4. Run app
```
go run .
```
