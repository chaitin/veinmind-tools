<h1 align="center"> veinmind-webshell </h1>

<p align="center">
veinmind-webshell 是由长亭科技自研的一款 Webshell 扫描工具 
</p>

## 功能特性

- 快速扫描镜像/容器中的 Webshell
- 支持`containerd`/`dockerd`容器运行时

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
chmod +x veinmind-webshell && ./veinmind-webshell scan xxx 
```
### 基于平行容器模式
确保机器上安装了`docker`以及`docker-compose`
#### Makefile 一键命令
```
make run.docker ARG="scan xxxx"
```
#### 自行构建镜像进行扫描
构建`veinmind-webshell`镜像
```
make build.docker
```
运行容器进行扫描
```
docker run --rm -it --mount 'type=bind,source=/,target=/host,readonly,bind-propagation=rslave' veinmind-webshell scan xxx
```


## 使用

1. 登录[百川平台](https://rivers.chaitin.cn/)，激活关山 Webshell 检测产品

![](../../../docs/veinmind-webshell/readme1.png)

2. 点击左下角组织配置创建 API Token (基础版每日限制检测 100 次， 高级版可联系问脉小助手/百川平台获取)

![](../../../docs/veinmind-webshell/readme2.png)

![](../../../docs/veinmind-webshell/readme3.png)

3. 使用token扫描指定镜像

```
./veinmind-webshell scan image [imageID/imageName] --token [关山token]
```
![](../../../docs/veinmind-webshell/scan_image_1.jpg)
4. 使用token扫描本地所有镜像
```
./veinmind-webshell scan image  --token [关山token]
```
![](../../../docs/veinmind-webshell/scan_image_2.jpg)

5.使用token扫描指定容器
```
./veinmind-webshell scan container [containerID/containerName] --token [关山token]
```
![](../../../docs/veinmind-webshell/scan_container_1.jpg)

5.使用token扫描本地所有容器
```
./veinmind-webshell scan container  --token [关山token]
```

![](../../../docs/veinmind-webshell/scan_container_2.jpg)

6. 指定输出格式
支持的输出格式：
- html
- json
- cli（默认）
```
./veinmind-webshell scan container [containerID/containerName] --token [token] -f html
```
生成的result.html效果如图：
![](../../../docs/veinmind-webshell/format.jpg)


