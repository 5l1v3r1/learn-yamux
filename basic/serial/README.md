basic serial example:  client with socket, server with serialport
=================================================================

# qemu parameter

```
-device virtserialport,chardev=charch0,id=channel0,name=agent.channel.0 
-chardev socket,id=charch0,path=/run/vc/vm/eed8b7781ae14bda827245ac053d56666d22d4087211212e08e17c404f9681b1/kata.sock,server,nowait
```

- server: the serial device name in windows is `\\.\agent.channel.0`
- client: use .sock file to connect server


# server

> run on win7 (qemu guest)

```
> cd server
> go build
> server.exe
Starting serial demo - server
open serial port \\.\agent.channel.0 okwait...
time="2019-08-02T18:52:34+08:00" level=info msg="===== begin to receive message from client ======"
time="2019-08-02T18:52:40+08:00" level=info msg="[1]received: hello"
time="2019-08-02T18:52:40+08:00" level=info msg="[2]received: world"
time="2019-08-02T18:52:40+08:00" level=info msg="[3]received: ping"
time="2019-08-02T18:52:40+08:00" level=info msg="sent: pong"
```


# client

> run on centos

```
$ cd client
$ go build
$ sudo ./client
  Starting serial demo - client
  INFO[0000] connect to /run/vc/vm/eed8b7781ae14bda827245ac053d56666d22d4087211212e08e17c404f9681b1/kata.sock
  INFO[0000] unix dial ok
  INFO[0000] ===== send message to server =====
  INFO[0000] [0]sent 'hello'
  INFO[0000] [1]sent 'world'
  INFO[0000] [2]sent 'ping'
  INFO[0000] ===== read message from server =====
  INFO[0000] received :pong
```