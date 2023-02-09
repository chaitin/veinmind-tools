
<h1 align="center"> veinmind-escalate </h1>

<p align="center">
veinmind-malicious 是由长亭科技自研的一款逃逸风险扫描工具 
</p>

## 功能特性

- 快速扫描容器中的逃逸风险 
- 支持 `docker`/`containerd` 容器运行时
- 支持`JSON`/`CSV`/`HTML`等多种报告格式输出

## 兼容性

- linux/amd64
- linux/386
- linux/arm64
- linux/arm

## 开始之前

### 安装方式一

请先安装`libveinmind`，安装方法可以参考[官方文档](https://github.com/chaitin/libveinmind)

确保机器上安装了`docker`以及`docker-compose`，并启动`ClamAV`。

```
chmod +x veinmind-escalate && ./veinmind-escalte extract && cd scripts && docker-compose pull && docker-compose up -d
```

如果您使用的是`VirusTotal`，则需要在环境变量或`scripts/.env`文件中声明`VT_API_KEY`
```
export VT_API_KEY=xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx
```

### 安装方式二

基于平行容器的模式，获取 `veinmind-escalate` 的镜像并启动
```
docker run --rm -it --mount 'type=bind,source=/,target=/host,readonly,bind-propagation=rslave' -v `pwd`:/tool/data veinmind/veinmind-escalate scan
```

或者使用项目提供的脚本启动
```
chmod +x parallel-container-run.sh && ./parallel-container-run.sh scan
```

## 使用

1.指定镜像名称或镜像ID并扫描 (需要本地存在对应的镜像)

```
./veinmind-escalate scan image [imageID/imageName]
```

2.扫描所有本地镜像

```
./veinmind-escalate scan image
```

3.指定容器名称或容器ID并扫描

```
./veinmind-escalate scan container [imageID/imageName]
```

4.扫描所有本地容器

```
./veinmind-escalate scan container
```

