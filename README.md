# 服务监控
Client端使用Golang语言编写，通过gRPC协议向Server端发送请求，上报服务器的CPU，内存，磁盘等监控信息。Server端的实现可以参照Google的gRPC官方文档。

server-agent使用cobra开发命令行工具，通过命令行工具可以启动Client端，配置Server端监听端口，以及Client端上报监控信息的时间间隔。


![image.jpg](image%2Fimage.jpg)


```
server-agent.exe -H 127.0.0.1 -p 9777 -f 3
```
可使用 -h 查看帮助
```
This is a tool that uses RPC (Remote Procedure Call) to push service information to the server.

Usage:
  server-agent [flags]

Flags:
  -f, --factor int    The factor value (default -1)
  -h, --help          help for server-agent
  -H, --host string   The host IP address
  -p, --port int      The port number (default -1)

```


启动成功后将看到以下日志：
```
2024/07/19 21:19:33 Host: 127.0.0.1, Port: 9777, Factor: 3
```


通过UPX压缩之后的文件大小只有2MB，非常方便部署。