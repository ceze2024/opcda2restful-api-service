# opc测试

## 安装LocalServiceComponents组件

## 安装
OPC_2.0_Core_Components-Setup

## 注册dll

- 一般注册x86

### 使用386编译的情形
把x32位的gbda_aut.dll文件拷贝到Windows的SysWOW64文件夹下


### 使用amd64编译的情形
把x64位的gbda_aut.dll文件拷贝到Windows的System32文件夹下

```cmd
regsvr32 gbda_aut.dll
```

环境变量设置，由于kepServer特殊，需要设置 386  ，打包需要设置为 amd64

```cmd 
set GOARCH=386
go build -o opcConnector-win32.exe

set GOARCH=amd64 
go build -o opcConnector-win64.exe
```


读取wincc时，需要保持wincc变量再画面中，否则数据无法更新

# 注册到windows服务
管理员cmd下执行：
```cmd 
opcConnector-win32 install 
```
 
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


 
