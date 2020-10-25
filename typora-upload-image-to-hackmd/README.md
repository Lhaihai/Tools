# typora 上传图片到 imgur

为了方便把本地图片上传到imgur，然后就可以直接复制文章所有内容到hackmd上了，不用再一个个手动找到图片文件再上传图片到hackmd

## 使用

![](https://i.imgur.com/uCDoxNu.png)

![](https://i.imgur.com/IU5h53E.png)

## 配置


### 获取imgur client-id

登录imgur账号，打开这个链接：https://api.imgur.com/oauth2/addclient

![](https://i.imgur.com/CvtzGlI.png)

完成后得到 client id


imgur api 上传接口示例：https://apidocs.imgur.com/?version=latest#c85c9dfc-7487-4de2-9ecd-66f727cf3139

![](https://i.imgur.com/3bKhXxf.png)


### 配置typora

打开typora，文件->偏好设置->图像

在上传服务设定里选 Custom Command，填入下列命令

```
python "F:\\path\\imgur_upload.py" client-id
```

![](https://i.imgur.com/VQ3pQX0.png)

点击验证图片上传，出现这样结果为正常

![](https://i.imgur.com/YfdfnIt.png)



