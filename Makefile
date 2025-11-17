migrate:
	migrate -path ./migrations -database 'postgres://admin:root123@0.0.0.0:5432/subscriptiondb?sslmode=disable' up
migrate-d:
	migrate -path ./migrations -database 'postgres://admin:root123@0.0.0.0:5432/subscriptiondb?sslmode=disable' down
migrate-force:
	migrate -path ./migrations -database 'postgres://admin:root123@0.0.0.0:5432/subscriptiondb?sslmode=disable' force 1
migrate-version:
	migrate -path ./migrations -database 'postgres://admin:root123@0.0.0.0:5432/subscriptiondb?sslmode=disable' version
swag:
	swag init -g cmd/main.go
build:
	docker-compose build
run:
	docker-compose up