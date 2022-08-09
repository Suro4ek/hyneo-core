proto auth:
	protoc --go_out=. \
                   --go-grpc_out=.  \
                   proto/auth.proto   
proto vk:
	protoc --go_out=. \
                   --go-grpc_out=. \
                   proto/vk.proto  