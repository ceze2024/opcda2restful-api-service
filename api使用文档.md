# api说明

## 读取变量

> 接口地址：

http://[host:port]/Read

> 接口方法：POST

请求header需要加参数：

Secret：ADSDWW1DSADSADSAWJJK （密钥为服务端配置）

> 请求body：

```json
{
    "tags":["GONGDUANNAME","AT_02F12R11/1.OosAct#Value"]
}
```

> 返回结果：

```json
{"code":0,"data":{"AT_02F12R11/1.OosAct#Value":{"Value":false,"Quality":64,"Timestamp":"2023-06-26T06:33:48Z"},"GONGDUANNAME":{"Value":"酶解","Quality":192,"Timestamp":"2023-06-22T21:46:30Z"}},"msg":"success"}
```

## 写入变量

> 接口地址：

http://[host:port]/Write

> 接口方法：POST

请求header需要加参数：

Secret：ADSDWW1DSADSADSAWJJK （密钥为服务端配置）

> 请求body：

```json
{
    "test.dev1.a": 2,
"tag2":1
}
```

> 返回结果：

```json
{
    "code": 0,
    "data": {},
    "msg": "success"
}
```



#### tags为wincc里的编码名称：

![](C:\Users\ThinkPad\AppData\Roaming\marktext\images\2023-06-26-14-44-51-image.png)
