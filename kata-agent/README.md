kata-agent client + server example
==================================

server: run in windows(guest vm)
client: run in linux(host)

yamux + grpc + virtio-serial

test grpc api: GetGuestDetails()

# server

## build
```
$ cd server
$ go build
```

## run server

```
c:\Users\admin> server.exe
time="2019-09-04T11:31:45+08:00" level=info msg=" s.initLogger()"
time="2019-09-04T11:31:45.51+08:00" level=info msg=announce debug_console=false device-handlers= name=kata-agent pid=306
4 source=agent storage-handlers= version=unknown
time="2019-09-04T11:31:45.511+08:00" level=info msg="setupDebugConsole: /dev/console"
time="2019-09-04T11:31:45.511+08:00" level=info msg="s.setupSignalHandler()"
time="2019-09-04T11:31:45.511+08:00" level=info msg="setupSignalHandler - begin"
time="2019-09-04T11:31:45.511+08:00" level=info msg="signalHandlerLoop - begin"
time="2019-09-04T11:31:45.512+08:00" level=info msg="close(errCh)"
time="2019-09-04T11:31:45.512+08:00" level=info msg="setupSignalHandler - end"
time="2019-09-04T11:31:45.512+08:00" level=info msg="s.initChannel()"
time="2019-09-04T11:31:45.512+08:00" level=info msg="newChannel - begin"
time="2019-09-04T11:31:45.513+08:00" level=info msg="newChannel:0"
time="2019-09-04T11:31:45.513+08:00" level=info msg="checkForSerialChannel: \\\\.\\agent.channel.0"
time="2019-09-04T11:31:45.515+08:00" level=info msg="newChannel - end"
time="2019-09-04T11:31:45.515+08:00" level=info msg="s.startGRPC()"
time="2019-09-04T11:31:45.516+08:00" level=info msg="s.waitForStopServer()"
time="2019-09-04T11:31:45.516+08:00" level=info msg="s.listenToUdevEvents()"
time="2019-09-04T11:31:45.516+08:00" level=info msg="s.wg.Wait()"
time="2019-09-04T11:31:45.516+08:00" level=info msg="agent grpc server starts" debug_console=false name=kata-agent pid=3
064 source=agent
time="2019-09-04T11:31:45.517+08:00" level=info msg="serialChannel.setup() - begin"
time="2019-09-04T11:31:45.517+08:00" level=info msg="setup() - open serial port \\\\.\\agent.channel.0 ok"
time="2019-09-04T11:31:45.518+08:00" level=info msg="serialChannel.setup() - end"
time="2019-09-04T11:31:45.518+08:00" level=info msg="serialChannel.wait() - begin"
time="2019-09-04T11:31:45.518+08:00" level=info msg="serialChannel.wait() - end"
time="2019-09-04T11:31:45.518+08:00" level=info msg="serialChannel.listen() - begin"
time="2019-09-04T11:31:45.518+08:00" level=info msg="listen() - init yamux server over serialport:\\\\.\\agent.channel.0
 ok"
time="2019-09-04T11:31:45.519+08:00" level=info msg="serialChannel.listen() - end"
time="2019-09-04T11:31:45.519+08:00" level=info msg="Waiting for stopServer signal..." debug_console=false name=kata-age
nt pid=3064 source=agent subsystem=stopserverwatcher

```

# client

## build
```
$ cd client
$ go build
```

## run client

```
$ sudo ./client
INFO[0002] version:0.0.1
```