basic example: remote execute command via socket(client) and serial(server)
===========================================================================


# server

> run on win7(qemu vm)

```
$ cd server
$ go build

$ server.exe
c:\Users\admin>server.exe
Starting remote_exec demo - server
open serial port \\.\agent.channel.0 ok
open console port \\.\console0 ok
wait...
time="2019-08-07T23:20:10+08:00" level=info msg="===== begin to receive message from client ======"
time="2019-08-07T23:20:53+08:00" level=info msg="[1]received: cd\n"
time="2019-08-07T23:20:53+08:00" level=info msg="execute cmd:[cmd /c cd\n]"
time="2019-08-07T23:20:53+08:00" level=info msg="sent result"
time="2019-08-07T23:20:56+08:00" level=info msg="[2]received: hostname\n"
time="2019-08-07T23:20:56+08:00" level=info msg="execute cmd:[cmd /c hostname\n]"
time="2019-08-07T23:20:56+08:00" level=info msg="sent result"
time="2019-08-07T23:21:12+08:00" level=info msg="[3]received: dir /w\n"
time="2019-08-07T23:21:12+08:00" level=info msg="execute cmd:[cmd /c dir /w\n]"
time="2019-08-07T23:21:12+08:00" level=info msg="sent result"
```


# client

run on centos

```
$ cd client
$ go build

$ sudo ./client
sudo ./client
Starting remote_exec demo - client
INFO[0000] unix dial /run/vc/vm/1cd65c2aefcb65ee2a2139373f4e041f35074b2d5a0f0c3f274ec2e9cdc18694/kata.sock ok
INFO[0000] unix dial /run/vc/vm/1cd65c2aefcb65ee2a2139373f4e041f35074b2d5a0f0c3f274ec2e9cdc18694/console.sock ok
c:\>cd
c:\Users\admin


c:\>hostname
VM-WIN7-TEST


c:\>dir /w
 驱动器 C 中的卷没有标签。
 卷的序列号是 D829-4657

 c:\Users\admin 的目录

[.]             [..]            .bashrc         .bash_profile   .minttyrc
.viminfo        client.exe      [Contacts]      [Desktop]       [Documents]
[Downloads]     [Favorites]     [go]            govendor.exe    [Links]
[Music]         [Pictures]      [Saved Games]   [Searches]      server.exe
[Videos]
               7 个文件     16,058,462 字节
              14 个目录 16,940,843,008 可用字节
```
