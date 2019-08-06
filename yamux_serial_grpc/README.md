# client

> run on centos

client + sock file

```
$ cd client
$ go build

$ sudo ./client
sudo ./client
INFO[0000]
INFO[0000] ===== send request 0 =====
INFO[0000] client is nil, create a new one
INFO[0000] New client                                    url=/run/vc/vm/1cd65c2aefcb65ee2a2139373f4e041f35074b2d5a0f0c3f274ec2e9cdc18694/kata.sock
INFO[0000] unix socket connected
INFO[0000] create yamux session
INFO[0000] create yamux stream
INFO[0002] grpc connection ok
INFO[0002] sent request: name:"world 0"
INFO[0002] receive response: Hello world 0
INFO[0002]
INFO[0002] ===== send request 1 =====
INFO[0002] use exist kata agent client
INFO[0002] sent request: name:"world 1"
INFO[0002] receive response: Hello world 1
INFO[0002]
INFO[0002] ===== send request 2 =====
INFO[0002] use exist kata agent client
INFO[0002] sent request: name:"world 2"
INFO[0002] receive response: Hello world 2
INFO[0002]
INFO[0002] done
INFO[0002] yamux session close
INFO[0002] yamux session close
```


# server

> run on win7(qemu vm)

server + serial port

```
> cd server
> go build

> server.exe
c:\Users\admin>server.exe
time="2019-08-06T21:51:40+08:00" level=info msg="connect to \\\\.\\Global\\agent.channel.0"
time="2019-08-06T21:51:40+08:00" level=info msg="start setup()"
time="2019-08-06T21:51:40+08:00" level=info msg="open serial port \\\\.\\Global\\agent.channel.0 ok"
time="2019-08-06T21:51:40+08:00" level=info msg="===== start listen() ====="
time="2019-08-06T21:51:40+08:00" level=info msg="init yamux server over serialport:\\\\.\\Global\\agent.channel.0 ok"
time="2019-08-06T21:51:40+08:00" level=info msg="start grpc server"
time="2019-08-06T21:51:49+08:00" level=info msg="handle SayHello: [name:\"world 0\" ]"
time="2019-08-06T21:51:49+08:00" level=info msg="handle SayHello: [name:\"world 1\" ]"
time="2019-08-06T21:51:50+08:00" level=info msg="handle SayHello: [name:\"world 2\" ]"
time="2019-08-06T21:51:50+08:00" level=warning msg="2019/08/06 21:51:50 [ERR] yamux: Failed to write header: Used to indicate that an operat
ion cannot continue without blocking for I/O.\n" component=yamux
time="2019-08-06T21:51:53+08:00" level=warning msg="2019/08/06 21:51:53 [ERR] yamux: Failed to read header: The handle is invalid.\n" compon
ent=yamux
time="2019-08-06T21:51:53+08:00" level=warning msg="grpc server exited with error:Used to indicate that an operation cannot continue without
 blocking for I/O."
time="2019-08-06T21:51:53+08:00" level=info msg="session close"
<wait for new session>
time="2019-08-06T21:51:53+08:00" level=info msg="open serial port \\\\.\\Global\\agent.channel.0 ok"
time="2019-08-06T21:51:53+08:00" level=info msg="===== start listen() ====="
time="2019-08-06T21:51:53+08:00" level=info msg="init yamux server over serialport:\\\\.\\Global\\agent.channel.0 ok"
time="2019-08-06T21:51:53+08:00" level=info msg="start grpc server"
```
