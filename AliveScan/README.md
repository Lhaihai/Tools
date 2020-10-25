
启发自klion师傅的内网存活自动化探测思路：https://mp.weixin.qq.com/s/AUgBlRjH_USaZXgmMDYzSg

因为用系统自带的ping太慢了，用golang实现了一下，默认使用512协程，一个B段5秒检测完，保存结果在alive.txt文件下

```
Usage:
    AliveScan host thread                 (默认使用512个并发协程)
    Example: ./AliveScan -host 192 -t 50  用50个并发协程探测192.168 B段存活的C段
             ./AliveScan -host 172        探测172.16-32 B段存活的C段
             ./AliveScan -host 10         探测10 A段存活的C段
             ./AliveScan -host 10.172     探测10.172 B段存活的C段
```
