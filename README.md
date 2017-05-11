# userClientIP
是用户端自检的小工具，可收集用户的本地IP，外网IP，traceroute信息

#用途
如果你的用户访问不了你的门户， 试试这个工具吧。




#文件介绍
1. getClientNatIPWindows.go      用于生成windowsX64平台下的代码
2. getClientNatIP.go             用户生成MacOS和LinuxX86_64的代码
3. getWebHostIP.go       提供一个简单的界面可以让用户下载各个平台的检测工具，并直接显示用户的ip和浏览器信息。


#生成文件 debug.txt 
```
内涵三个信息：
1   客户端的本地ip
2   客户端访问ip.cn ，myexternalip.com ， ipip.net  三个平台看到的IP
3   从客户端trace网站的信息
```


## 开发环境
golang 1.8 

## 作者介绍
yihongfei  QQ:413999317   MAIL:yihf@liepin.com

CCIE 38649


## 寄语
为网络自动化运维尽绵薄之力
