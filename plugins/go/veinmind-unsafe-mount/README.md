<h1 align="center"> veinmind-unsafe-mount </h1>

<p align="center">
veinmind-unsafe-mount 是由长亭科技自研的一款容器不安全挂载目录扫描工具 
</p>

## 功能特性

- 快速扫描容器中的 unsafe-mount
- 支持`containerd`/`dockerd`容器运行时

## 兼容性

- linux/amd64
- linux/386
- linux/arm64

## 开始之前

### 安装方式一

请先安装`libveinmind`，安装方法可以参考[官方文档](https://github.com/chaitin/libveinmind)

### 安装方式二

基于平行容器的模式，获取 `veinmind-unsafe-mount` 的镜像并启动
```
docker run --rm -it --mount 'type=bind,source=/,target=/host,readonly,bind-propagation=rslave' veinmind/veinmind-unsafe-mount scan-container
```

或者使用项目提供的脚本启动
```
chmod +x parallel-container-run.sh && ./parallel-container-run.sh scan
```