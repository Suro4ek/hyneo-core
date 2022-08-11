proto auth:
	protoc --go_out=. -I=proto \
				   --go_opt=paths=import \
				   --go-grpc_out=. -I=proto \
				   --go-grpc_opt=paths=import \
				   proto/auth.proto
proto service:
	protoc --go_out=. -I=proto \
				   --go_opt=paths=import \
				   --go-grpc_out=. -I=proto \
				   --go-grpc_opt=paths=import \
				   proto/service.proto
go linux:
	GOARCH=amd64 GOOS=linux go build cmd/main.go cmd/services.go

go build:
	go build cmd/main.go cmd/services.go