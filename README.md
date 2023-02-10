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
<p align="center"> veinmind, 中文名为<b>问脉</b>，寓意 <b>容器安全见筋脉，望闻问切治病害。</b> 旨在成为云原生领域的一剂良方 </p>
</p>
<p align="center"> 中文文档 | <a href="README.en.md">English</a> </p>

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
./run.sh scan image
```


## 🔨 工具列表

| 工具                                                        | 功能                | 
|-----------------------------------------------------------|-------------------|
| [veinmind-runner](veinmind-runner/README.md)              | 扫描工具运行宿主          |
| [veinmind-malicious](plugins/go/veinmind-malicious)       | 扫描容器/镜像中的恶意文件     |
| [veinmind-weakpass](plugins/go/veinmind-weakpass)         | 扫描容器/镜像中的弱口令      |
| [veinmind-log4j2](plugins/go/veinmind-log4j2)             | 扫描容器/镜像中的log4j2漏洞 |
| [veinmind-sensitive](plugins/python/veinmind-sensitive)   | 扫描镜像中的敏感信息        |
| [veinmind-backdoor](plugins/python/veinmind-backdoor)     | 扫描镜像中的后门          |
| [veinmind-history](plugins/python/veinmind-history)       | 扫描镜像中的异常历史命令      |
| [veinmind-vuln](plugins/go/veinmind-vuln)                 | 扫描容器/镜像中的资产信息和漏洞  |
| [veinmind-webshell](plugins/go/veinmind-webshell)         | 扫描镜像中的 Webshell   |
| [veinmind-unsafe-mount](plugins/go/veinmind-unsafe-mount) | 扫描容器中的不安全挂载目录     |
| [veinmind-iac](plugins/go/veinmind-iac)                   | 扫描IaC文件           |
| [veinmind-escalate](plugins/go/veinmind-escalate)         | 扫描容器/镜像中的逃逸风险     |
    
PS: 目前所有工具均已支持平行容器的方式运行

## 🧑‍💻 编写插件

可以通过 example 快速创建一个 veinmind-tools 插件, 具体查看 [veinmind-example](example/)  

## ☁️ 云原生设施兼容性
| 名称                                                          | 类别    | 是否兼容 |
|-------------------------------------------------------------|-------|------|
| [Jenkins](https://github.com/chaitin/veinmind-jenkins)      | CI/CD | ✔️   |
| Gitlab CI                                                   | CI/CD | ✔️   |
| [Github Action](https://github.com/chaitin/veinmind-action) | CI/CD | ✔️   |
| DockerHub                                                   | 镜像仓库  | ✔️   |
| Docker Registry                                             | 镜像仓库  | ✔️   |
| Harbor                                                      | 镜像仓库  | ✔️   |
| Docker                                                      | 容器运行时 | ✔️   |
| Containerd                                                  | 容器运行时 | ✔️   |
| Kubernetes                                                  | 集群    | ✔️   |

## 🛴 工作原理
![](docs/architecture.png)

## 🏘️ 联系我们
1. 您可以通过 GitHub Issue 直接进行 Bug 反馈和功能建议。
2. 扫描下方二维码可以通过添加问脉小助手，以加入问脉用户讨论群进行详细讨论

![](docs/veinmind-group-qrcode.jpg)

## ✨ 404星链计划
<img src="https://github.com/knownsec/404StarLink-Project/raw/master/logo.png" width="30%">

veinmind-tools 现已加入 [404星链计划](https://github.com/knownsec/404StarLink)
