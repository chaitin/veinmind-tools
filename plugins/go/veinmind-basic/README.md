<h1 align="center"> veinmind-basic </h1>

<p align="center">
veinmind-basic 是由长亭科技自研的一款镜像/容器详细信息扫描工具 
</p>

## 功能特性

- 快速扫描镜像/容器中的详细信息
- 支持 `docker`/`containerd` 容器运行时
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
chmod +x veinmind-basic && ./veinmind-basic scan xxx 
```
### 基于平行容器模式
确保机器上安装了`docker`以及`docker-compose`
#### Makefile 一键命令
```
make run.docker ARG="scan xxxx"
```
#### 自行构建镜像进行扫描
构建`veinmind-basic`镜像
```
make build.docker
```
运行容器进行扫描
```
docker run --rm -it --mount 'type=bind,source=/,target=/host,readonly,bind-propagation=rslave' veinmind-basic scan xxx
```

## 使用参数

1.指定镜像名称或镜像ID并扫描 (需要本地存在对应的镜像)

```
./veinmind-basic scan image [imagename/imageid]
```
![](../../../docs/veinmind-basic/basic_scan_image_1.jpeg)

2.扫描所有本地镜像

```
./veinmind-basic scan image
```
![](../../../docs/veinmind-basic/basic_scan_image_2.jpeg)

3.指定容器名称或容器ID并扫描 (需要本地存在对应的容器)
```
./veinmind-basic scan container [containerName/containerid]
```
![](../../../docs/veinmind-basic/basic_scan_container_1.jpeg)

4.扫描所有本地容器
```
./veinmind-basic scan container
```
![](../../../docs/veinmind-basic/basic_scan_container_2.jpeg)

5.指定输出类型
  支持的输出格式：
- html
- json
- cli（默认）
```
./veinmind-basic scan image [imageID/imageName] -f html
```
生成的result.html效果如图：
![](../../../docs/veinmind-basic/basic_format_1.jpg)
