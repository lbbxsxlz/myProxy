## myProxy
	使用Go语言实现[socks5协议](https://github.com/lbbxsxlz/myProxy/blob/master/SOCKS5_RFC1928_en.md)

	经验证亦可以在arm架构的设备上代理网络

## 编译
	GOOS=linux GOARCH=arm go build myProxy.go

## 使用
	e.p. ./myProxy --listen 0.0.0.0:9999
