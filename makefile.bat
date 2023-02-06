protoc --proto_path=./proto --go-grpc_out=. service.proto
protoc --proto_path=./proto --go_out=. service.proto
pause