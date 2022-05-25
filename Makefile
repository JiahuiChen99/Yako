gen_proto:
	protoc --proto_path=src/grpc/proto/ --go-grpc_opt=require_unimplemented_servers=false --go_out=src/grpc/ --go-grpc_out=src/grpc/ src/grpc/proto/*.proto

clean:
	rm src/grpc/pb/*.go

run_master:
	./src/yako_master/YakoMaster $(ip) $(port) $(zk_ip) $(zk_port) $(mqtt_ip) $(mqtt_port)

build_master:
	go build src/yako_master/YakoMaster.go

run_agent:
	./src/yako_node/YakoAgent $(ip) $(port) $(zk_ip) $(zk_port)

build_agent:
	go build src/yako_node/YakoAgent.go

run_agent_iot:
	./src/yako_agent_iot/YakoAgentIoT $(ip) $(port) $(mqtt_ip) $(mqtt_port)

build_agent_iot:
	go build src/yako_agent_iot/YakoAgentIoT.go