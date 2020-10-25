把C段IP和端口输出为ip列表，方便导入其它工具使用



```
Usage of c2ip.exe:
  -f string
        -f domain or ip or cip list
  -h string
        host
  -o string
        -o url.txt (default "url.txt")
  -p string
        Enter the Ports!! (default "80")
```



常用命令

```
c2ip.exe -f ip.txt -p 80,8080-8090 -o url.txt
```



ip.txt

```
10.1.1.1/24
10.1.2.1
www.baidu.com
```

