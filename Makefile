gen:
	protoc --proto_path=src/grpc/proto/ --go-grpc_opt=require_unimplemented_servers=false --go_out=src/grpc/ --go-grpc_out=src/grpc/ src/grpc/proto/*.proto

clean:
	rm src/grpc/pb/*.go

run:
	go run src/main.go

build:
	go build src/main.gop