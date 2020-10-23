## myProxy
使用Go语言实现[socks5协议](https://github.com/lbbxsxlz/myProxy/blob/master/SOCKS5_RFC1928_en.md)
代码中附上了socks5的协议细节

在x86\amd64\arm\arm64架构的设备上代理网络验证成功

## 编译
```
go build
```
其他系统与架构请指定对应的操作系统与架构，例如
```
GOOS=linux GOARCH=arm go build myProxy.go
GOOS=linux GOARCH=amd64 go build myProxy.go
GOOS=windows GOARCH=386 go build myProxy.go
```

## 使用
e.p. ./myProxy --listen 0.0.0.0:9999

同时配合浏览器配置代理服务
