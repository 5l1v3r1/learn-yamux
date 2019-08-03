
# how to use gRPC
- Define a service in a .proto file.
- Generate server and client code using the protocol buffer compiler.
- Use the Go gRPC API to write a simple client and server for your service.


# define proto

helloworld/helloworld.proto



# generate code


## install proto

```
$ go get -u -v github.com/golang/protobuf/protoc-gen-go

$ sudo wget https://github.com/protocolbuffers/protobuf/releases/download/v3.7.1/protoc-3.7.1-osx-x86_64.zip -P /usr/local
$ cd /usr/local
$ sudo unzip protoc-3.7.1-osx-x86_64.zip


```

## generate go code

```
$ cd pkg/grpc
$ protoc -I protos/ protos/helloworld.proto --go_out=plugins=grpc:./protos
$ ll protos/helloworld.pb.go
-rw-r--r--  1 xjimmy  staff   4.3K Apr 13 10:58 protos/helloworld.pb.go
```

# write server code

server/main.go


# write client code

client/main.go


# run

## start gRPC server

> run on macos 

```
$ go get -v github.com/golang/protobuf/proto
$ go get -v google.golang.org/grpc
$ go get -v github.com/sirupsen/logrus
$ go get -v golang.org/x/net

$ go build

$ ./server
INFO[0000] Start gRPC Server on port ::50051            
INFO[0010] [172.16.87.134:49178]client connected, received request: [jimmy], sent response [hello jimmy] 
```

## start gRPC client

> run on win7(qemu guest)

```
$ go build

$ ./client 172.16.87.1:50051 xjimmy
time="2019-08-04T00:16:44+08:00" level=info msg="connect to 172.16.87.1:50051"
time="2019-08-04T00:16:44+08:00" level=info msg="sent request: [name=jimmy], received response: [hello jimmy]"
```
