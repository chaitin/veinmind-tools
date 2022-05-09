<p align="center">
  <img src="https://dinfinite.oss-cn-beijing.aliyuncs.com/image/20220428154824.png" width="120">
</p>
<h1 align="center"> veinmind-tools </h1>
<p align="center">
  <a href="https://veinmind.chaitin.com/docs/">Documentation</a> 
</p>

<p align="center">
<img src="https://img.shields.io/github/v/release/chaitin/veinmind-tools.svg" />
<img src="https://img.shields.io/github/release-date/chaitin/veinmind-tools.svg?color=blue&label=update" />
<img src="https://img.shields.io/badge/go report-A+-brightgreen.svg" />

<p align="center"> veinmind-tools 是由长亭科技自研，基于 <a href="https://github.com/chaitin/libveinmind">veinmind-sdk</a> 打造的容器安全工具集 </p>
</p>

## 🔥 Demo
![](https://dinfinite.oss-cn-beijing.aliyuncs.com/image/20220415144819.gif)


## 🕹️ 快速开始
### 1. 确保机器上正确安装 docker
```
docker info
```
### 2. 安装 [veinmind-runner](https://github.com/chaitin/veinmind-tools/tree/master/veinmind-runner) 镜像
```
docker pull veinmind/veinmind-runner:latest
```
### 3. 下载 [veinmind-runner](https://github.com/chaitin/veinmind-tools/tree/master/veinmind-runner) 平行容器启动脚本
```
wget -q https://download.veinmind.tech/scripts/veinmind-runner-parallel-container-run.sh -O run.sh && chmod +x run.sh
```
### 4. 快速扫描本地镜像
```
./run.sh scan-host
```


## 🔨 工具列表

|  工具 | 功能  | 
|---|---|
|  [veinmind-runner](https://github.com/chaitin/veinmind-tools/tree/master/veinmind-runner) | 扫描工具运行宿主 |
|  [veinmind-malicious](https://github.com/chaitin/veinmind-tools/tree/master/veinmind-malicious) | 扫描镜像中的恶意文件  |
|  [veinmind-weakpass](https://github.com/chaitin/veinmind-tools/tree/master/veinmind-weakpass)  | 扫描镜像中的弱口令  |
|  [veinmind-sensitive](https://github.com/chaitin/veinmind-tools/tree/master/veinmind-sensitive) | 扫描镜像中的敏感信息  |
|  [veinmind-backdoor](https://github.com/chaitin/veinmind-tools/tree/master/veinmind-backdoor) | 扫描镜像中的后门 |
|  [veinmind-history](https://github.com/chaitin/veinmind-tools/tree/master/veinmind-history) | 扫描镜像中的异常历史命令 |
|  [veinmind-asset](https://github.com/chaitin/veinmind-tools/tree/master/veinmind-asset) | 扫描镜像中的资产信息 |
    
PS: 目前所有工具均已支持平行容器的方式运行

## 🏘️ 联系我们
1. 您可以通过 GitHub Issue 直接进行 Bug 反馈和功能建议。
2. 扫描下方二维码可以通过添加问脉小助手，以加入问脉用户讨论群进行详细讨论

![](docs/veinmind-group-qrcode.jpg)

## ✨ 404星链计划
<img src="https://github.com/knownsec/404StarLink-Project/raw/master/logo.png" width="30%">

veinmind-tools 现已加入 [404星链计划](https://github.com/knownsec/404StarLink)
