<h1 align="center"> veinmind-trace </h1>

<p align="center">
veinmind-trace 是由长亭科技自研的一款容器安全检测工具
</p>

## 功能特性

+ 快速扫描容器中的异常进程:
    1. 隐藏进程(mount -o bind方式)
    2. 反弹shell的进程
    3. 带有挖矿、黑客工具、可疑进程名的进程
    4. 包含 Ptrace 的进程
+ 快速扫描容器中的异常文件系统:
    1. 敏感目录权限异常
    2. cdk 工具利用痕迹检测
+ 快速扫描容器中的异常用户:
    1. uid=0 的非root账户
    2. uid相同的用户
+ 支持`containerd`/`dockerd`容器运行时

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
chmod +x veinmind-trace && ./veinmind-trace scan xxx 
```

### 基于平行容器模式

确保机器上安装了`docker`以及`docker-compose`

#### Makefile 一键命令

```
make run.docker ARG="scan xxxx"
```

#### 自行构建镜像进行扫描

构建`veinmind-trace`镜像

```
make build.docker
```

运行容器进行扫描

```
docker run --rm -it --mount 'type=bind,source=/,target=/host,readonly,bind-propagation=rslave' veinmind-trace scan xxx
```

## 使用

1. 扫描本地所有容器

```
./veinmind-trace scan container
```
