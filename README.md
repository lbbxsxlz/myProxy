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

同时在需要上网的机器上使用浏览器配置代理服务

以谷歌 chrome 浏览器为例：<br>
１、在 google 浏览器中添加 SwitchyOmega 插件<br>
在Chrome地址栏输入**chrome://extensions**即可打开扩展程序，然后拖动后缀名为.crx的**SwitchyOmega**安装文件到扩展程序中进行安装，片刻即可安装完成。<br>
Chrome较高版本已经不支持以上描述的.crx安装，需要将后缀crx改成zip，然后解压zip压缩包。然后在扩展程序页面中，找到右上角的开发者模式，将开发者模式打开，<br>
点击“加载已解压的扩展程序”铵钮，打开文件选择窗口，选择刚才解压的目录，然后点击“选择文件夹”按钮即可。<br>
２、点击浏览器的右上角的扩展程序按钮，选择“proxy”命令，进入“SwitchyOmega 页面”。<br>
３、在 SwitchyOmega 页面，选择代理协议为“SOCKS5”，输入代理服务器地址和端口（运行myproxy的IP与端口），填完之后，点击“应用选项”按钮。<br>
４、点击浏览器的右上角的扩展程序按钮，选择“proxy”命令。<br>
５、在浏览器中输入想访问的域名或IP即可，Please　enjoy！
