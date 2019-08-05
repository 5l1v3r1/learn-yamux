# client

> run on centos

client + sock file

```
$ cd client
$ go build

$ sudo ./client
INFO[0000] connect to /run/vc/vm/1cd65c2aefcb65ee2a2139373f4e041f35074b2d5a0f0c3f274ec2e9cdc18694/kata.sock
INFO[0000]
INFO[0000] ===== send request 0 =====
INFO[0000] start sendReq
INFO[0000] start connect
INFO[0000] client is nil, create a new one
INFO[0000] New client                                    url=/run/vc/vm/1cd65c2aefcb65ee2a2139373f4e041f35074b2d5a0f0c3f274ec2e9cdc18694/kata.sock
INFO[0000] NewAgentClient, kataURL: /run/vc/vm/1cd65c2aefcb65ee2a2139373f4e041f35074b2d5a0f0c3f274ec2e9cdc18694/kata.sock
INFO[0000] grpcAddr:unix:////run/vc/vm/1cd65c2aefcb65ee2a2139373f4e041f35074b2d5a0f0c3f274ec2e9cdc18694/kata.sock parsedAddr:/run/vc/vm/1cd65c2aefcb65ee2a2139373f4e041f35074b2d5a0f0c3f274ec2e9cdc18694/kata.sock
INFO[0000] start agentDialer()
INFO[0000] return yamux dialer
INFO[0000] before grpc.DialContext: grpcAddr:unix:////run/vc/vm/1cd65c2aefcb65ee2a2139373f4e041f35074b2d5a0f0c3f274ec2e9cdc18694/kata.sock
INFO[0000] start yamux dialer, sock:unix:////run/vc/vm/1cd65c2aefcb65ee2a2139373f4e041f35074b2d5a0f0c3f274ec2e9cdc18694/kata.sock timeout:19.999876557s
INFO[0000] start unixDialer()
INFO[0000] waiting...
INFO[0000] start net.DialTimeout sock:////run/vc/vm/1cd65c2aefcb65ee2a2139373f4e041f35074b2d5a0f0c3f274ec2e9cdc18694/kata.sock timeout:19.999876557s
INFO[0000] commonDialer conn ok
INFO[0000] receive conn ok
INFO[0000] create yamux client
INFO[0000] start create yamux stream
INFO[0000] yamux create stream ok
INFO[0000] after grpc.DialContext: grpcAddr:unix:////run/vc/vm/1cd65c2aefcb65ee2a2139373f4e041f35074b2d5a0f0c3f274ec2e9cdc18694/kata.sock, error:<nil>
INFO[0000] grpc.DialContext ok
INFO[0000] start installReqFunc
INFO[0000] connect ok
INFO[0000] reqHandlers:map[grpc.SayHello:0x81a730]
INFO[0000] request:name:"world"
INFO[0000] get handler for grpc.SayHello
INFO[0000] call handler
INFO[0000] response:Hello world
INFO[0000]
INFO[0000] ===== send request 1 =====
INFO[0000] start sendReq
INFO[0000] start connect
INFO[0000] use exist kata agent client
INFO[0000] connect ok
INFO[0000] reqHandlers:map[grpc.SayHello:0x81a730]
INFO[0000] request:name:"world"
INFO[0000] get handler for grpc.SayHello
INFO[0000] call handler
INFO[0000] response:Hello world
INFO[0000]
INFO[0000] ===== send request 2 =====
INFO[0000] start sendReq
INFO[0000] start connect
INFO[0000] use exist kata agent client
INFO[0000] connect ok
INFO[0000] reqHandlers:map[grpc.SayHello:0x81a730]
INFO[0000] request:name:"world"
INFO[0000] get handler for grpc.SayHello
INFO[0000] call handler
INFO[0000] response:Hello world
INFO[0000] done
```


# server

> run on win7(qemu vm)

server + serial port

```
> cd server
> go build

> server.exe
time="2019-08-05T18:50:35+08:00" level=info msg="connect to \\\\.\\Global\\agent.channel.0"
time="2019-08-05T18:50:35+08:00" level=info msg="start setup()"
time="2019-08-05T18:50:35+08:00" level=info msg="open serial port \\\\.\\Global\\agent.channel.0 ok"
time="2019-08-05T18:50:35+08:00" level=info msg="===== start listen() ====="
time="2019-08-05T18:50:35+08:00" level=info msg="init yamux server over serialport:\\\\.\\Global\\agent.channel.0 ok"
time="2019-08-05T18:50:35+08:00" level=info msg="return makeUnaryInterceptor()"
time="2019-08-05T18:50:35+08:00" level=info msg="start grpc server"
time="2019-08-05T18:50:35+08:00" level=warning msg="2019/08/05 18:50:35 wait...\n" component=yamux
time="2019-08-05T18:50:38+08:00" level=info msg="start makeUnaryInterceptor()"
time="2019-08-05T18:50:38+08:00" level=warning msg="Creating gRPC context as none found"
time="2019-08-05T18:50:38+08:00" level=error msg="failed to get peer of client"
time="2019-08-05T18:50:38+08:00" level=info msg="receive gRPC request: [name:\"world\" ] [<nil>]"
time="2019-08-05T18:50:38+08:00" level=info msg="start makeUnaryInterceptor()"
time="2019-08-05T18:50:38+08:00" level=warning msg="Creating gRPC context as none found"
time="2019-08-05T18:50:38+08:00" level=error msg="failed to get peer of client"
time="2019-08-05T18:50:38+08:00" level=info msg="receive gRPC request: [name:\"world\" ] [<nil>]"
time="2019-08-05T18:50:38+08:00" level=info msg="start makeUnaryInterceptor()"
time="2019-08-05T18:50:38+08:00" level=warning msg="Creating gRPC context as none found"
time="2019-08-05T18:50:38+08:00" level=error msg="failed to get peer of client"
time="2019-08-05T18:50:38+08:00" level=info msg="receive gRPC request: [name:\"world\" ] [<nil>]"
```
