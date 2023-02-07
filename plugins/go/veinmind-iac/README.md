<h1 align="center"> veinmind-iac </h1>

<p align="center">
veinmind-iac 用于扫描IaC(Infrastructure as Code) 文件内的风险问题
</p>

## 功能特性

- 支持 `dockerfile/kubernetes` IaC类型文件
- 支持指定目录自动递归扫描

## 兼容性

- linux/amd64
- linux/386
- linux/arm64
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
chmod +x veinmind-iac && ./veinmind-iac scan xxx 
```
### 基于平行容器模式
确保机器上安装了`docker`以及`docker-compose`
#### Makefile 一键命令
```
make run.docker ARG="scan xxxx"
```
#### 自行构建镜像进行扫描
构建`veinmind-iac`镜像
```
make build.docker
```
运行容器进行扫描
```
docker run --rm -it --mount 'type=bind,source=/,target=/host,readonly,bind-propagation=rslave' veinmind-iac scan xxx
```

## 使用参数

1. 指定扫描IaC文件

```
./veinmind-iac scan iac IACFILE
```
![img.png](../../../docs/veinmind-iac/iac_scan_iac_01.jpg)


2. 指定扫描目录下可能存在的IaC文件类型

```
./veinmind-iac scan iac PATH
```
![img.png](../../../docs/veinmind-iac/iac_scan_iac_02.jpg)

3. 指定扫描特定的IaC文件类型

```
./veinmind-iac scan iac --iac-type kubernetes/dockerfile IACFILE/PATH
```
![img.png](../../../docs/veinmind-iac/iac_scan_iac_03.jpg)

4. 指定输出格式 
支持的输出格式：
- html
- json
- cli（默认）
```
./veinmind-iac scan iac -f html IACFILE/PATH
```
生成的result.html效果如图：
![img.png](../../../docs/veinmind-iac/iac_scan_iac_04.jpg)
