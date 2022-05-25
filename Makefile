gen_proto:
	protoc --proto_path=src/grpc/proto/ --go-grpc_opt=require_unimplemented_servers=false --go_out=src/grpc/ --go-grpc_out=src/grpc/ src/grpc/proto/*.proto

clean:
	rm src/grpc/pb/*.go

run_master:
	./src/yako_master/YakoMaster $(ip) $(port) $(zk_ip) $(zk_port) $(mqtt_ip) $(mqtt_port)

build_master:
	go build src/yako_master/YakoMaster.go

run_agent:
	go run src/yako_node/YakoAgent.go

build_agent:
	go build src/yako_node/YakoAgent.go