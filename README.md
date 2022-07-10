# Yako
Yako an orchestrator that determines the viability and handles the deployment of computing services into a multi-layered distributed heterogeneous environment with support for traditional and IoT computing nodes.
The project is powered by Golang, Apache Zookeeper, gRPC and MQTT.

YakoUI is the client to interact with the Yako platform.

This project is part of the [Heterogeneous Computing Farm](https://github.com/JiahuiChen99/Heterogeneous-Computing-Farm), my Bachelor of Science in Computer Science senior thesis.
For more in depth explanation please refer to the [senior thesis paper](https://github.com/JiahuiChen99/Heterogeneous-Computing-Farm/blob/main/Heterogeneous%20Computing%20Farm.pdf).

## 🧰 Prerequisites

- [Golang v1.17.x](https://go.dev/)
- [Make](https://www.gnu.org/software/make/manual/make.html)
- [Apache Zookeeper](https://zookeeper.apache.org/)
- [gRPC](https://grpc.io/)
- [Protocol buffers (Protobuf)](https://developers.google.com/protocol-buffers)

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
Follow https://grpc.io/docs/languages/go/quickstart/
