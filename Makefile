include .env

run:
	go run cmd/api/main.go http

build:
	GOARCH=amd64 GOOS=linux go build -o main_ydesetiawan94 cmd/api/main.go

deploy:
	scp -i w1key main_ydesetiawan94 ubuntu@52.221.209.87:~

migration_setup:
	psql -U ${DB_USERNAME} -c "CREATE DATABASE ${DB_NAME};"

migration_up:
	migrate -path db/migrations/ -database "postgresql://${DB_USERNAME}:${DB_PASSWORD}@${DB_HOST}:${DB_PORT}/${DB_NAME}?sslmode=disable" -verbose up

migration_down:
	@read -p "Enter VERSION: " VERSION; \
	migrate -path db/migrations -database "postgresql://${DB_USERNAME}:${DB_PASSWORD}@${DB_HOST}:${DB_PORT}/${DB_NAME}?sslmode=disable" -verbose down $$VERSION

migration_fix:
	@read -p "Enter VERSION: " VERSION; \
	migrate -path db/migrations -database "postgresql://${DB_USERNAME}:@${DB_HOST}:${DB_PORT}/${DB_NAME}?sslmode=disable" force $$VERSION

migration_up_showcase:
	migrate -database "postgres://postgres:iatuyachie1Hae4Maih5izee1vie6Ooxu@projectsprint-db.cavsdeuj9ixh.ap-southeast-1.rds.amazonaws.com:5432/postgres?sslrootcert=ap-southeast-1-bundle.pem&sslmode=verify-full" -path db/migrations up
