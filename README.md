# Yako
Yako an orchestrator that determines the viability and handles the deployment of computing services into a multi-layered distributed heterogeneous environment with support for traditional and IoT computing nodes.
The project is powered by Golang, Apache Zookeeper, gRPC and MQTT.

YakoUI is the client to interact with the Yako platform.

This project is part of the [Heterogeneous Computing Farm](https://github.com/JiahuiChen99/Heterogeneous-Computing-Farm), my Bachelor of Science in Computer Science senior thesis.
For more in depth explanation please refer to the [senior thesis paper](https://github.com/JiahuiChen99/Heterogeneous-Computing-Farm/blob/main/Heterogeneous%20Computing%20Farm.pdf).

## 🧰 Prerequisites

- [Golang v1.17.x](https://go.dev/) or higher
- [Make](https://www.gnu.org/software/make/manual/make.html)
- [Apache Zookeeper](https://zookeeper.apache.org/)
- [gRPC](https://grpc.io/)
- [Protocol buffers (Protobuf)](https://developers.google.com/protocol-buffers)
- [Mosquitto (MQTT Broker)](https://mosquitto.org/)

## ⚙ Installation
Make sure that the correct go version is installed in your system by running `go version`.
The project provides a Makefile with all the directives to either build or run both YakoMaster & YakoAgent.

### gRPC & Protocol Buffers
gRPC related RPC procedures must be generated before proceeding with the project setup.
Generate the Go gRPC source code by executing `make gen_proto`. 
This will take all **.proto** files from **src/grpc/proto** and create all the boilerplate in **src/grpc/yako**.

Make sure to install these Go plugins used for protocol buffers compilation. Do not run `sudo apt install protobuf-compiler`.

```shell
go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
export PATH="$PATH:$(go env GOPATH)/bin"
```

Further information can be retrieved from the official Go gRPC quickstart guide
Follow https://grpc.io/docs/languages/go/quickstart/.


Download and install golang project dependencies with `go mod tidy`. And run the following 3 Makefile rules to build the software.

```shell
# Builds YakoMaster
make build_master

# Builds YakoMaster
make build_agent

# Builds YakoAgent (IoT)
make build_agent_iot
````

## 🚀️ Execution

Steps to run the Yako platform:

1. [Run Apache Zookeeper](https://github.com/JiahuiChen99/Yako#)
2. [Run mosquitto MQTT Broker](https://github.com/JiahuiChen99/Yako#)
3. [Run YakoMaster](https://github.com/JiahuiChen99/Yako#)
4. [Run YakoAgent](https://github.com/JiahuiChen99/Yako#)
5. [Run YakoAgent (IoT)](https://github.com/JiahuiChen99/Yako#)

### Service Registry
Yako uses a Service Registry to keep track of the active available nodes within the cluster. This functionality is provided by Apache Zookeeper.
A configuration file is provided in **src/utils/zookeeper/zoo.cnf**.


To run or stop zookeeper, run with the following commands:

```shell
# If ZK is registered to the path
zkServer start
zkServer stop

# Otherwise go to the folder where ZK is installed
./zkServer start
./zkServer stop

# ZK also provides a CLI tool to interact with
./zkCli
```

### MQTT Broker
Yako platform uses MQTT protocol that provides a publish/subscribe pattern to interact with YakoAgent (IoT) devices. The broker used for this purpose Mosquitto
A configuration file is located in **src/utils/mqtt/mosquitto.conf**.

Run the following command to start mosquitto broker with the configuration file.

```shell
mosquitto -c mosquitto.conf
```

### YakoMaster
The orchestrator itself accepts a total of 6 arguments.

| Argument  | Description                |
| --------- | ---------------            |
| ip        | device IP                  |
| port      | service port               |
| zk_ip     | Zookeeper IP               |
| zk_port   | Zookeeper port             |
| mqtt_ip   | mosquitto MQTT broker IP   |
| mqtt_port | mosquitto MQTT broker port |

```shell
# Makefile rule
make run_master ip=<IP> port=<Port> zk_ip=<ZK IP> zk_port=<ZK Port> mqtt_ip=<MQTT IP> mqtt_port=<MQTT Port>

# Manual
/YakoMaster <IP> <Port> <ZK IP> <ZK Port> <MQTT IP> <MQTT Port>
```

### YakoAgent
To run the agent for computing nodes, four arguments must be assigned:

| Argument | Description    |
| -------- | -------------- |
| ip       | device IP      |
| port     | service port   |
| zk_ip    | Zookeeper IP   |
| zk_port  | Zookeeper port |

```shell
# Makefile rule
make run_agent ip=<IP> port=<Port> zk_ip=<ZK IP> zk_port=<ZK Port>

# Manual
/YakoAgent <IP> <Port> <ZK IP> <ZK Port>
```