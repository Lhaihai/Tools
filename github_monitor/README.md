

起因是在群里看到有人发了个github监测工具

https://github.com/M4tir/Github-Monitor

但是是用易语言写的，不是很懂易语言，所以就亲自操刀用python乱写一通

因为之前写过一点钉钉机器人，就不用邮箱了，改为用钉钉机器人发送通知



使用方法：

代码第16行的token修改为钉钉机器人token，在20行修改secret，如何设置钉钉机器人请自行谷歌

```
python3 github_monitor.py
```

每五分钟监测一下github是否有新出的CVE或者RCE



![](https://i.imgur.com/tayWTaz.png)

test.db 保存了2019和2020的CVE项目，一些POC和RCE项目

![](https://i.imgur.com/kgYjQON.png)

log.txt 保存新发现的项目
