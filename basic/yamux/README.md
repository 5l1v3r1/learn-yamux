basic yamux demo: one session， multiple stream
===============================================

# server

> run on macos

```
$ cd server
$ go build

//启动服务
$ ./server
[main]Starting yamux demo
[main]wait done here
INFO[0000] start listener accept on 0.0.0.0:4444        
INFO[0006] 
== new client connected, creating yamux server session now ==  clientID="172.16.87.134:45204"
INFO[0006] [0]hello[1] [2]world                          clientID="172.16.87.134:45204" streamID=7
INFO[0006] [0]hello[1] [2]world                          clientID="172.16.87.134:45204" streamID=3
INFO[0006] [0]hello[1] [2]world                          clientID="172.16.87.134:45204" streamID=5
INFO[0006] [0]hello[1] [2]world                          clientID="172.16.87.134:45204" streamID=1
INFO[0006] session closed                                clientID="172.16.87.134:45204"
```


# client

> run on win7 (qemu guest)

```
//启动
$ cd client
$ go build

$ ./client 172.16.87.1:4444
Starting yamux demo - client
INFO[0000] connect to 172.16.87.1:4444
INFO[0000] creating client session
INFO[0000] create yamux session ok
INFO[0000] waiting...
INFO[0000] opening stream 7
INFO[0000] opening stream 1
INFO[0000] opening stream 5
INFO[0000] opening stream 3
INFO[0000] done
```