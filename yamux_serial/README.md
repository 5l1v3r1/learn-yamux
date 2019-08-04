# client

> run on centos

client + sock file

```
$ cd client
$ go build

$ sudo ./client
INFO[0000] connect to /run/vc/vm/1cd65c2aefcb65ee2a2139373f4e041f35074b2d5a0f0c3f274ec2e9cdc18694/kata.sock
INFO[0000] unix dial ok
INFO[0000] create yamux client ok
INFO[0000] send ping_0
INFO[0001] receive response: [pong_0]
INFO[0002] send ping_1
INFO[0002] receive response: [pong_1]
INFO[0003] send ping_2
INFO[0003] receive response: [pong_2]
INFO[0004] exit normal
INFO[0004] exit with error:<nil>
```


# server

> run on win7(qemu vm)

server + serial port

```
> cd server
> go build

> server.exe
time="2019-08-05T10:31:07+08:00" level=info msg="connect to \\\\.\\Global\\agent.channel.0"
time="2019-08-05T10:31:07+08:00" level=info msg="start setup()"
time="2019-08-05T10:31:07+08:00" level=info msg="open serial port \\\\.\\Global\\agent.channel.0 ok"
time="2019-08-05T10:31:07+08:00" level=info msg="===== start listen() ====="
time="2019-08-05T10:31:07+08:00" level=info msg="init yamux server over serialport:\\\\.\\Global\\agent.channel.0 ok"
time="2019-08-05T10:31:07+08:00" level=info msg="session Accept"
time="2019-08-05T10:31:07+08:00" level=warning msg="2019/08/05 10:31:07 wait...\n" component=yamux
time="2019-08-05T10:31:10+08:00" level=warning msg="2019/08/05 10:31:10 wait...\n" component=yamux
time="2019-08-05T10:31:13+08:00" level=warning msg="2019/08/05 10:31:13 wait...\n" component=yamux
time="2019-08-05T10:31:17+08:00" level=info msg=Recv
time="2019-08-05T10:31:17+08:00" level=info msg="recv [ping_0]"
time="2019-08-05T10:31:17+08:00" level=info msg="send response [pong]"
time="2019-08-05T10:31:18+08:00" level=info msg="recv [ping_1]"
time="2019-08-05T10:31:18+08:00" level=info msg="send response [pong]"
time="2019-08-05T10:31:19+08:00" level=info msg="recv [ping_2]"
time="2019-08-05T10:31:19+08:00" level=info msg="send response [pong]"
time="2019-08-05T10:31:20+08:00" level=warning msg="2019/08/05 10:31:20 wait...\n" component=yamux
time="2019-08-05T10:31:20+08:00" level=error msg="stream read EOF"
time="2019-08-05T10:31:20+08:00" level=info msg="close stream"
time="2019-08-05T10:31:20+08:00" level=warning msg="2019/08/05 10:31:20 [ERR] yamux: Failed to write header: Used to ind
icate that an operation cannot continue without blocking for I/O.\n" component=yamux
time="2019-08-05T10:31:20+08:00" level=info msg="session is over."
time="2019-08-05T10:31:20+08:00" level=info msg="session close"
<wait for new session>
time="2019-08-05T10:31:23+08:00" level=warning msg="2019/08/05 10:31:23 [ERR] yamux: Failed to read header: The handle i
s invalid.\n" component=yamux
time="2019-08-05T10:31:23+08:00" level=info msg="open serial port \\\\.\\Global\\agent.channel.0 ok"
time="2019-08-05T10:31:23+08:00" level=info msg="===== start listen() ====="
time="2019-08-05T10:31:23+08:00" level=info msg="init yamux server over serialport:\\\\.\\Global\\agent.channel.0 ok"
time="2019-08-05T10:31:23+08:00" level=info msg="session Accept"
time="2019-08-05T10:31:23+08:00" level=warning msg="2019/08/05 10:31:23 wait...\n" component=yamux
time="2019-08-05T10:31:26+08:00" level=warning msg="2019/08/05 10:31:26 wait...\n" component=yamux
```
