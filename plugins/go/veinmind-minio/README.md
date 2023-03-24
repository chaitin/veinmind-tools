<h1 align="center"> veinmind-minio </h1>

<p align="center">
veinmind-minio  主要用于扫描 CVE-2023-28432 漏洞专项检查。
</p>

## 功能特性

- 快速扫描容器/镜像中的minio CVE-2023-28432风险
- 支持`JSON`/`CLI`/`HTML`等多种报告格式输出

## 兼容性

- linux/amd64
- linux/386
- linux/arm64
- linux/arm

## 使用方式

### BinaryStandalone(Beta)
静态编译版本，无需`libveinmind`依赖，下载即可用。

从[下载仓库](https://download.veinmind.tech/plugins%2Fveinmind-minio)对应系统版本的二进制文件。

执行:
`chmod +x veinmind-minio-amd64 && ./veinmind-minio-amd64 scan image`

更多详情，可查看：https://stack.chaitin.com/techblog/detail?id=62

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
chmod +x veinmind-minio && ./veinmind-minio scan xxx 
```
### 基于平行容器模式
确保机器上安装了`docker`以及`docker-compose`
#### Makefile 一键命令
```
make run.docker ARG="scan xxxx"
```
#### 自行构建镜像进行扫描
构建`veinmind-minio`镜像
```
make build.docker
```
运行容器进行扫描
```
docker run --rm -it --mount 'type=bind,source=/,target=/host,readonly,bind-propagation=rslave' registry.veinmind.tech/veinmind/veinmind-minio scan xxx
```

## 使用参数

1.指定镜像名称或镜像ID并扫描 (需要本地存在对应的镜像)

```
./veinmind-minio scan image [imageID/imageName]
```
![](../../../docs/veinmind-minio/minio_scan_1.png)
2.扫描所有本地镜像

```
./veinmind-minio scan image
```

3.指定容器名称或容器ID并扫描

```
./veinmind-minio scan container [containerID/containerName]
```
![](../../../docs/veinmind-minio/minio_scan_2.png)
4.扫描所有本地容器

```
./veinmind-minio scan container
```

5.指定输出格式
支持的输出格式：
- html
- json
- cli（默认）
```
./veinmind-minio scan container [containerID/containerName] -f html
```
生成的result.html效果如图：

![](../../../docs/veinmind-minio/minio_scan_3.png)

