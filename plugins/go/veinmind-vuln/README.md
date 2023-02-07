<h1 align="center"> veinmind-vuln </h1>

<p align="center">
veinmind-vuln 用于扫描容器/镜像的资产和漏洞信息</p>

## 功能特性

- 扫描镜像/容器的OS信息
- 扫描镜像/容器内系统安装的packages
- 扫描镜像/容器内应用安装的libraries
- 扫描镜像/容器是否存在已知cve (beta)

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
chmod +x veinmind-vuln && ./veinmind-vuln scan xxx 
```
### 基于平行容器模式
确保机器上安装了`docker`以及`docker-compose`
#### Makefile 一键命令
```
make run.docker ARG="scan xxxx"
```
#### 自行构建镜像进行扫描
构建`veinmind-vuln`镜像
```
make build.docker
```
运行容器进行扫描
```
docker run --rm -it --mount 'type=bind,source=/,target=/host,readonly,bind-propagation=rslave' veinmind-vuln scan xxx
```

## 使用参数

1.指定镜像名称或镜像ID并扫描 (需要本地存在对应的镜像)

```
./veinmind-vuln scan image [imageID/imageName]
```
![](../../../docs/veinmind-vuln/vuln_scan_image_01.jpg)
2.扫描所有本地镜像

```
./veinmind-vuln scan image
```
![](../../../docs/veinmind-vuln/vuln_scan_image_02.jpg)
3.指定容器名称或容器ID并扫描

```
./veinmind-vuln scan container [containerID/containerName]
```
![](../../../docs/veinmind-vuln/vuln_scan_container_01.jpg)


4.扫描所有本地容器

```
./veinmind-vuln scan container
```
![](../../../docs/veinmind-vuln/vuln_scan_container_02.jpg)


5.指定输出格式
支持的输出格式：
- html
- json
- cli（默认）
```
./veinmind-vuln scan image [imageID/imageName] -f html
```
生成的result.html效果如图：

![](../../../docs/veinmind-vuln/vuln_scan_image_05.jpg)
6.显示详细信息
```
./veinmind-vuln scan image [imageID/imageName] -v
```
![](../../../docs/veinmind-vuln/vuln_scan_image_06.jpg)
7.显示特定类型的信息
```
./veinmind-vuln scan image [imageID/imageName] --type [os/python/npm/jar.....]
```
![](../../../docs/veinmind-vuln/vuln_scan_image_07.jpeg)


8.仅扫描资产信息
```
./veinmind-vuln scan image [imageID/imageName] --only-asset
```
![](../../../docs/veinmind-vuln/vuln_scan_image_08.jpg)
