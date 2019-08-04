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
INFO[0000] send ping
```


# server

> run on win7(qemu vm)

server + serial port

```
> cd server
> go build

> server.exe
time="2019-08-04T12:07:10+08:00" level=info msg="connect to \\\\.\\Global\\agent.channel.0"
time="2019-08-04T12:07:10+08:00" level=info msg="start setup()"
time="2019-08-04T12:07:10+08:00" level=info msg="open serial port \\\\.\\Global\\agent.channel.0 ok"
time="2019-08-04T12:07:10+08:00" level=info msg="start listen()"
time="2019-08-04T12:07:10+08:00" level=info msg="init yamux server over serialport:\\\\.\\Global\\agent.channel.0 ok"
time="2019-08-04T12:07:10+08:00" level=info msg="session Accept"
time="2019-08-04T12:07:10+08:00" level=warning msg="2019/08/04 12:07:10 wait...\n" component=yamux
time="2019-08-04T12:07:13+08:00" level=warning msg="2019/08/04 12:07:13 wait...\n" component=yamux
time="2019-08-04T12:07:16+08:00" level=info msg=Recv
time="2019-08-04T12:07:16+08:00" level=info msg="Recv: [ID=1] ping"
time="2019-08-04T12:07:16+08:00" level=info msg="Recv: [ID=1] ping"
time="2019-08-04T12:07:17+08:00" level=info msg="Recv: [ID=1] ping"
time="2019-08-04T12:07:18+08:00" level=info msg="Recv: [ID=1] ping"
time="2019-08-04T12:07:19+08:00" level=info msg="Recv: [ID=1] ping"
<connect againt>
time="2019-08-04T12:07:25+08:00" level=warning msg="2019/08/04 12:07:25 [ERR] yamux: duplicate stream declared\n" compon
ent=yamux
time="2019-08-04T12:07:25+08:00" level=error msg="stop old stream"
time="2019-08-04T12:07:25+08:00" level=warning msg="2019/08/04 12:07:25 [ERR] yamux: Failed to write header: The handle
is invalid.\n" component=yamux
time="2019-08-04T12:07:25+08:00" level=info msg="session over."
```
