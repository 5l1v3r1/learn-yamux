remote execute command via socket(client) and serial(server)
============================================================


# server

> run on win7(qemu vm)

```
$ cd server
$ go build

$ server.exe
Starting exec_serial demo - server
[read] open serial port \\.\agent.channel.0 ok
[write] open serial port \\.\console0 ok
wait...
time="2019-08-08T12:43:44+08:00" level=info msg="===== begin to receive message from client ======"
time="2019-08-08T12:44:30+08:00" level=info msg="[1]received: dir\n"
time="2019-08-08T12:44:30+08:00" level=info msg="execute cmd:[cmd /c dir\n]"
time="2019-08-08T12:44:30+08:00" level=info msg="finish execute cmd:[cmd /c dir\n]"
time="2019-08-08T12:44:39+08:00" level=info msg="[2]received: ipconfig\n"
time="2019-08-08T12:44:39+08:00" level=info msg="execute cmd:[cmd /c ipconfig\n]"
time="2019-08-08T12:44:39+08:00" level=info msg="finish execute cmd:[cmd /c ipconfig\n]"
time="2019-08-08T12:44:43+08:00" level=info msg="[3]received: ping 172.19.0.242\n"
time="2019-08-08T12:44:43+08:00" level=info msg="execute cmd:[cmd /c ping 172.19.0.242\n]"
time="2019-08-08T12:44:46+08:00" level=info msg="finish execute cmd:[cmd /c ping 172.19.0.242\n]"
```


# client

run on centos

```
$ cd client
$ go build

$ sudo ./client
Starting exec_serial demo - client
INFO[0000] [write] unix dial /run/vc/vm/1cd65c2aefcb65ee2a2139373f4e041f35074b2d5a0f0c3f274ec2e9cdc18694/kata.sock ok
INFO[0000] [read] unix dial /run/vc/vm/1cd65c2aefcb65ee2a2139373f4e041f35074b2d5a0f0c3f274ec2e9cdc18694/console.sock ok
C:\Users\admin>dir
 驱动器 C 中的卷没有标签。
 卷的序列号是 D829-4657

 c:\Users\admin 的目录

2019/08/07  23:21    <DIR>          .
2019/08/07  23:21    <DIR>          ..
2019/08/03  23:09                54 .bashrc
2019/08/07  19:31                99 .bash_profile
2019/08/03  23:08                41 .minttyrc
2019/08/07  23:21            22,734 .viminfo
2019/08/04  00:16        11,690,496 client.exe
2019/06/27  18:47    <DIR>          Contacts
2019/08/03  23:07    <DIR>          Desktop
2019/06/27  18:47    <DIR>          Documents
2019/07/10  11:10    <DIR>          Downloads
2019/06/27  18:47    <DIR>          Favorites
2019/08/05  07:09    <DIR>          go
2019/08/05  07:15         1,703,936 govendor.exe
2019/06/27  18:47    <DIR>          Links
2019/06/27  18:47    <DIR>          Music
2019/06/27  18:47    <DIR>          Pictures
2019/06/27  18:47    <DIR>          Saved Games
2019/06/27  18:47    <DIR>          Searches
2019/08/08  12:40         2,633,216 server.exe
2019/06/27  18:47    <DIR>          Videos
               7 个文件     16,050,576 字节
              14 个目录 16,938,360,832 可用字节
C:\Users\admin>ipconfig

Windows IP 配置


以太网适配器 本地连接 2:

   连接特定的 DNS 后缀 . . . . . . . :
   本地链接 IPv6 地址. . . . . . . . : fe80::4982:6f5a:78c:5f61%13
   IPv4 地址 . . . . . . . . . . . . : 172.19.0.242
   子网掩码  . . . . . . . . . . . . : 255.255.255.0
   默认网关. . . . . . . . . . . . . : 172.19.0.1

隧道适配器 isatap.{3973F837-CC7C-41E9-B7C5-E063B60F3E69}:

   媒体状态  . . . . . . . . . . . . : 媒体已断开
   连接特定的 DNS 后缀 . . . . . . . :
C:\Users\admin>ping 172.19.0.242

正在 Ping 172.19.0.242 具有 32 字节的数据:
来自 172.19.0.242 的回复: 字节=32 时间<1ms TTL=128
来自 172.19.0.242 的回复: 字节=32 时间<1ms TTL=128
来自 172.19.0.242 的回复: 字节=32 时间<1ms TTL=128
来自 172.19.0.242 的回复: 字节=32 时间<1ms TTL=128

172.19.0.242 的 Ping 统计信息:
    数据包: 已发送 = 4，已接收 = 4，丢失 = 0 (0% 丢失)，
往返行程的估计时间(以毫秒为单位):
    最短 = 0ms，最长 = 0ms，平均 = 0ms
C:\Users\admin>
```
