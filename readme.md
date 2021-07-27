# GoService

将普通 exe 包装成一个 Windows 服务。理论上跨平台，但是只测试过 Windows

基于[service](https://github.com/kardianos/service)

## 用法

### 编写配置文件

goservive.yml

```yaml
# 服务名称
servicename: GoHelloWeb
# 要执行的exe所在的路径
basedir: "D:\\goservice\\test"
# 要执行的exe文件
bin: "hello.exe"
# 要传给exe的参数
args:
  - -p
  - 127.0.0.1:8088
```

### 安装服务

**要在有管理员权限的命令行下执行**

```
goservice.exe -config goservive.yml -action install
```

然后就可以在 windows 的服务面板上看到该服务了

### 卸载服务

```
goservice.exe -config goservive.yml -action uninstall
```
