<h1 align="center"> veinmind-webshell </h1>

<p align="center">
veinmind-webshell 是由长亭科技自研的一款镜像 Webshell 扫描工具 
</p>

## 功能特性

- 快速扫描镜像中的 Webshell
- 支持`containerd`/`dockerd`容器运行时

## 兼容性

- linux/amd64
- linux/386
- linux/arm64

## 开始之前

### 安装方式一

请先安装`libveinmind`，安装方法可以参考[官方文档](https://github.com/chaitin/libveinmind)

### 安装方式二

基于平行容器的模式，获取 `veinmind-webshell` 的镜像并启动
```
docker run --rm -it --mount 'type=bind,source=/,target=/host,readonly,bind-propagation=rslave' veinmind/veinmind-webshell scan --token [关山token]
```

或者使用项目提供的脚本启动
```
chmod +x parallel-container-run.sh && ./parallel-container-run.sh scan --token [关山token]
```

## 使用

1. 登录[百川平台](https://rivers.chaitin.cn/)，激活关山 Webshell 检测产品
   ![](../../../docs/veinmind-webshell/readme1.png)

2. 点击左下角组织配置创建 API Token (基础版每日限制检测 100 次， 高级版可联系问脉小助手/百川平台获取)
   ![](../../../docs/veinmind-webshell/readme2.png)
   ![](../../../docs/veinmind-webshell/readme3.png)
3. 执行 `veinmind-webshell` 时填入创建的 `token`
```
./veinmind-webshell scan --token [关山token]
```