<h1 align="center"> veinmind-weakpass </h1>

<p align="center">
veinmind-weakpass 是由长亭科技自研的一款镜像弱口令扫描工具 
</p>

## 功能特性

- 快速扫描 镜像/容器 中的弱口令
- 支持弱口令宏定义
- 支持并发扫描弱口令
- 支持自定义用户名以及字典
- 支持`containerd`/`dockerd`容器运行时

## 兼容性

- linux/amd64
- linux/386
- linux/arm64
- linux/arm

## 开始之前

### 安装方式一

请先安装`libveinmind`，安装方法可以参考[官方文档](https://github.com/chaitin/libveinmind)

### 安装方式二

基于平行容器的模式，获取 `veinmind-weakpass` 的镜像并启动
```
docker run --rm -it --mount 'type=bind,source=/,target=/host,readonly,bind-propagation=rslave' veinmind/veinmind-weakpass scan image
```

或者使用项目提供的脚本启动
```
chmod +x parallel-container-run.sh && ./parallel-container-run.sh scan image
```

## 使用

1.指定镜像名称或镜像ID并扫描 (需要本地存在对应的镜像)

```
./veinmind-weakpass scan image [imagename/imageid]
```

2.指定容器名称或容器ID并扫描 (需要本地存在对应的容器)

```
./veinmind-weakpass scan container [containername/containerid]
```

3.扫描所有本地镜像

```
./veinmind-weakpass scan image
```

4.指定容器运行时类型
```
./veinmind-weakpass scan image --containerd
```

容器运行时类型
- dockerd
- containerd

5.指定扫描用户名类型
```
./veinmind-weakpass scan image -u username
```

6.指定自定义扫描字典
```
./veinmind-weakpass scan image -d ./pass.dict
```

7.指定自定义扫描的服务
```
./veinmind-weakpass scan image -a ssh,mysql,redis
```
- 目前已经支持的服务

    | serverName | version |
    |:----------:|:-------:|
    |     ssh    |   all   |
    |    mysql   |   8.X   |
    |    redis   |   all   |
    |   tomcat   |   all   |

8.解压默认字典到本地磁盘
```
./veinmind-weakpass extract
```

## 演示
1.扫描指定镜像名称 `test` 的所有服务
![](../../../docs/veinmind-weakpass/weakpasscandemo1.png)
2.扫描指定镜像名称 `test` 的 `ssh` 服务
![](../../../docs/veinmind-weakpass/weakpasscandemo2.png)
2.扫描所有镜像的 `ssh` 服务
![](../../../docs/veinmind-weakpass/weakpasscandemo3.png)
