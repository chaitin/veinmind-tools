<h1 align="center"> veinmind-log4j2 </h1>

<p align="center">
veinmind-log4j2 主要用于扫描存在CVE-2021-44228 漏洞的 log4j jar 文件。
</p>

## 功能特性

- 快速扫描容器/镜像中的log4j2风险
- 支持 fat jar、jar 中包含 jar 等情况的检测
- 支持`JSON`/`CLI`/`HTML`等多种报告格式输出

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
chmod +x veinmind-log4j2 && ./veinmind-log4j2 scan xxx 
```
### 基于平行容器模式
确保机器上安装了`docker`以及`docker-compose`
#### Makefile 一键命令
```
make run.docker ARG="scan xxxx"
```
#### 自行构建镜像进行扫描
构建`veinmind-log4j2`镜像
```
make build.docker
```
运行容器进行扫描
```
docker run --rm -it --mount 'type=bind,source=/,target=/host,readonly,bind-propagation=rslave' veinmind-log4j2 scan xxx
```

## 使用参数

1.指定镜像名称或镜像ID并扫描 (需要本地存在对应的镜像)

```
./veinmind-log4j2 scan image [imageID/imageName]
```
![](../../../docs/veinmind-log4j2/log4j2_scan_image_1.jpg)
2.扫描所有本地镜像

```
./veinmind-log4j2 scan image
```
![](../../../docs/veinmind-log4j2/log4j2_scan_image_2.jpg)
3.指定容器名称或容器ID并扫描

```
./veinmind-log4j2 scan container [containerID/containerName]
```
![](../../../docs/veinmind-log4j2/log4j2_scan_container_1.jpg)


4.扫描所有本地容器

```
./veinmind-log4j2 scan container
```
![](../../../docs/veinmind-log4j2/log4j2_scan_container_2.jpg)

5.指定输出格式
支持的输出格式：
- html
- json
- cli（默认）
```
./veinmind-log4j2 scan container [containerID/containerName] -f html
```
生成的result.html效果如图：

![](../../../docs/veinmind-log4j2/log4j2_format.jpg)
