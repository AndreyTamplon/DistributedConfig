# ==============================================================================
run:
	go run ./cmd/config_service/main.go
build:
	go build ./cmd/config_service/main.go
test:
	go test ./...

# ==============================================================================
# migrate postgres
DB_NAME = dc
DB_USER =
DB_PASS =
DB_HOST = localhost
force:
	migrate -path migrations -database "postgres://$(DB_HOST)/$(DB_NAME)?sslmode=disable&user=$(DB_USER)&password=$(DB_PASS)" force 1
version:
	migrate -path migrations -database "postgres://$(DB_HOST)/$(DB_NAME)?sslmode=disable&user=$(DB_USER)&password=$(DB_PASS)" version
migrate_up:
	migrate -path migrations -database "postgres://$(DB_HOST)/$(DB_NAME)?sslmode=disable&user=$(DB_USER)&password=$(DB_PASS)" up
migrate_down:
	migrate -path migrations -database "postgres://$(DB_HOST)/$(DB_NAME)?sslmode=disable&user=$(DB_USER)&password=$(DB_PASS)" down
# ==============================================================================
# proto
proto:
	protoc --proto_path=internal/delivery/proto --go_out=internal/delivery/ internal/delivery/proto/*.proto
	protoc --proto_path=internal/delivery/proto --go-grpc_out=require_unimplemented_servers=false:internal/delivery/ internal/delivery/proto/*.proto
gateway:
	protoc --grpc-gateway_out=internal/delivery/proto --grpc-gateway_opt=logtostderr=true --grpc-gateway_opt=paths=source_relative --proto_path=internal/delivery/proto internal/delivery/proto/*.proto