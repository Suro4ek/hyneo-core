proto auth:
	protoc --go_out=. \
                   --go-grpc_out=.  \
                   proto/auth.proto   
proto service:
	protoc --go_out=. \
                   --go-grpc_out=. \
                   proto/service.proto  