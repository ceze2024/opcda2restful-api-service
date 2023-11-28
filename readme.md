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

