<h1 align="center"> veinmind-weakpass </h1>

<p align="center">
veinmind-weakpass 是由长亭科技自研的一款容器/镜像弱口令扫描工具 
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
## 使用方式

### 基于可执行文件

请先安装`libveinmind`，安装方法可以参考[官方文档](https://github.com/chaitin/libveinmind)
#### Makefile 一键命令

```
make run ARG="scan xxx"
```
#### 自行编译可执行文件进行扫描

编译可执行文件
```
make build
```
运行可执行文件进行扫描
```
chmod +x veinmind-weakpass && ./veinmind-weakpass scan xxx 
```
### 基于平行容器模式
确保机器上安装了`docker`以及`docker-compose`
#### Makefile 一键命令
```
make run.docker ARG="scan xxxx"
```
#### 自行构建镜像进行扫描
构建`veinmind-weakpass`镜像
```
make build.docker
```
运行容器进行扫描
```
docker run --rm -it --mount 'type=bind,source=/,target=/host,readonly,bind-propagation=rslave' veinmind-weakpass scan xxx
```

## 使用参数
1.指定镜像名称或镜像ID并扫描 (需要本地存在对应的镜像)
```
./veinmind-weakpass scan image [imagename/imageid]
```
![](../../../docs/veinmind-weakpass/weakpass_scan_image_1.jpeg)
2.指定容器名称或容器ID并扫描 (需要本地存在对应的容器)
```
./veinmind-weakpass scan container [containername/containerid]
```
![](../../../docs/veinmind-weakpass/weakpass_scan_container_1.jpg)
3.扫描所有本地镜像
```
./veinmind-weakpass scan image 
```
![](../../../docs/veinmind-weakpass/weakpass_scan_image_3.jpeg)

4.扫描本地所有容器
```
./veinmind-weakpass scan container 
```
![](../../../docs/veinmind-weakpass/weakpass_scan_container_2.jpg)

5.指定扫描用户名类型
```
./veinmind-weakpass scan image -u username
```
![](../../../docs/veinmind-weakpass/weakpass_scan_image_5.jpeg)

6.指定自定义扫描字典
```
./veinmind-weakpass scan image -d ./pass.dict
```
![](../../../docs/veinmind-weakpass/weakpass_scan_image_6.jpeg)

7.指定自定义扫描的服务
```
./veinmind-weakpass scan image -s ssh,mysql,redis
```
目前已经支持的服务

| serverName | version |
|:----------:|:-------:|
|     ssh    |   all   |
|    mysql   |   8.X   |
|    redis   |   all   |
|   tomcat   |   all   |
|     ftp    |   all   |
![](../../../docs/veinmind-weakpass/weakpass_scan_image_7.jpeg)

8.解压默认字典到本地磁盘
```
./veinmind-weakpass extract
```
9.指定输出格式
支持的输出格式：
- html
- json
- cli（默认）
```
./veinmind-weakpass scan image [imageID/imageName] -f html
```
生成的result.html效果如图：
![](../../../docs/veinmind-weakpass/weakpass_scan_image_9.jpg)
