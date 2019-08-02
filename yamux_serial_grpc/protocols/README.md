
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
$ protoc -I grpc/ grpc/helloworld.proto --go_out=plugins=grpc:./grpc
$ ll grpc/helloworld.pb.go
-rw-r--r--  1 xjimmy  staff   4.3K Apr 13 10:58 grpc/helloworld.pb.go
```

# write server code

server/main.go


# write client code

client/main.go


# run

## start gRPC server
```
$ go get -v github.com/golang/protobuf/proto
$ go get -v google.golang.org/grpc
$ go get -v github.com/sirupsen/logrus
$ go get -v golang.org/x/net

$ go build
$ ./server
Start gRPC Server on port ::50051
INFO[0002] receive gRPC request: world  

```

## start gRPC client
```
$ go build
$ ./client
2019/04/13 13:23:35 Receiver gRPC response: Hello world
```