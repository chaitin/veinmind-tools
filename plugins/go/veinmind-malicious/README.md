<h1 align="center"> veinmind-malicious </h1>

<p align="center">
veinmind-malicious 是由长亭科技自研的一款镜像恶意文件扫描工具 
</p>

## 功能特性

- 快速扫描镜像中的恶意文件 (目前支持`ClamAV`以及`VirusTotal`)
- 支持 `docker`/`containerd` 容器运行时
- 支持`JSON`/`CLI`/`HTML`等多种报告格式输出

## 兼容性

- linux/amd64
- linux/386
- linux/arm64
- linux/arm

## 使用方式

开始之前请先确保机器上安装了clamav，并设置配置文件

```
cp dockerfiles/clamd.conf /etc/clamav/clamd.conf
```

如果您使用的是`VirusTotal`，则需要在环境变量或`scripts/.env`文件中声明`VT_API_KEY`
```
export VT_API_KEY=xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx
```
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
chmod +x veinmind-malicious && ./veinmind-malicious scan xxx 
```
### 基于平行容器模式
确保机器上安装了`docker`以及`docker-compose`
#### Makefile 一键命令
```
make run.docker ARG="scan xxxx"
```
#### 自行构建镜像进行扫描
构建`veinmind-malicious`镜像
```
make build.docker
```
运行容器进行扫描
```
docker run --rm -it --mount 'type=bind,source=/,target=/host,readonly,bind-propagation=rslave' veinmind-malicious scan xxx
```

## 使用参数

1.指定镜像名称或镜像ID并扫描 (需要本地存在对应的镜像)

```
./veinmind-malicious scan image [imagename/imageid]
```
![](../../../docs/veinmind-malicious/malicious_scan_image1.jpg)
2.扫描所有本地镜像

```
./veinmind-malicious scan image
```
![](../../../docs/veinmind-malicious/malicious_scan_image2-1.jpg)

![](../../../docs/veinmind-malicious/malicious_scan_image2-2.jpg)



3.指定输出报告格式
支持的输出格式：
- html
- json
- cli（默认）

```
./veinmind-malicious scan image -f html
```
生成的result.html效果如图：
![](../../../docs/veinmind-malicious/malicious_format.jpg)
