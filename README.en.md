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

<p align="center"> veinmind-tools is self-developed by <a href="https://www.chaitin.cn/en/"> chaitin technology </a>，a container security toolset based on <a href="https://github.com/chaitin/libveinmind">veinmind-sdk</a>  </p>
</p>

## 🔥 Demo

![](https://dinfinite.oss-cn-beijing.aliyuncs.com/image/20220415144819.gif)

## 🕹️ Quick Start

### 1. Make sure docker is installed correctly on the machine

```
docker info
```

### 2. Install [veinmind-runner](https://github.com/chaitin/veinmind-tools/tree/master/veinmind-runner) image

```
docker pull veinmind/veinmind-runner:latest
```

### 3. Download [veinmind-runner](https://github.com/chaitin/veinmind-tools/tree/master/veinmind-runner) parallel container startup script

```
wget -q https://download.veinmind.tech/scripts/veinmind-runner-parallel-container-run.sh -O run.sh && chmod +x run.sh
```

### 4. Quick scan local images

```
./run.sh scan image
```

## 🔨 Toolset

| Tool                                                                 | Description                                            | 
|----------------------------------------------------------------------|--------------------------------------------------------|
| [veinmind-runner](veinmind-runner/README.en.md)                      | scanner host                                           |
| [veinmind-malicious](plugins/go/veinmind-malicious/README.en.md)     | Scan containers/images for malicious files             |
| [veinmind-weakpass](plugins/go/veinmind-weakpass/README.en.md)       | scan containers/images for weak passwords              |
| [veinmind-log4j2](plugins/go/veinmind-log4j2/README.en.md)           | scan containers/images for log4j2(CVE-2021-44228)      |
| [veinmind-sensitive](plugins/python/veinmind-sensitive/README.en.md) | scan images for sensitive information                  |
| [veinmind-backdoor](plugins/python/veinmind-backdoor/README.en.md)   | scan images for backdoors                              |
| [veinmind-history](plugins/python/veinmind-history/README.en.md)     | scan images for abnormal history commands              |
| [veinmind-vuln](plugins/go/veinmind-vuln/README.en.md)               | scan containers/images for asset information and vulns |
| [veinmind-webshell](plugins/go/veinmind-webshell)                    | scan containers/images for webshell                    |
| [veinmind-unsafe-mount](plugins/go/veinmind-unsafe-mount)            | scan containers for unsafe mount                       |
| [veinmind-iac](plugins/go/veinmind-iac)                              | scan IaC file                                          |
| [veinmind-escalate](plugins/go/veinmind-escalate)                    | Scan containers/images for escape risks                |

PS: All tools currently support running in parallel containers

## 🧑‍💻 Coding Plugins

Use exmaple to create a veinmind-tool plugin quickly, see more at [veinmind-example](example/README.en.md)

## ☁️ Cloud-native infrastructure compatibility

| Name                                                        | Type     | Compatibility |
|-------------------------------------------------------------|----------|---------------|
| [Jenkins](https://github.com/chaitin/veinmind-jenkins)      | CI/CD    | ✔️            |
| Gitlab CI                                                   | CI/CD    | ✔️            |
| [Github Action](https://github.com/chaitin/veinmind-action) | CI/CD    | ✔️            |
| DockerHub                                                   | Registry | ✔️            |
| Docker Registry                                             | Registry | ✔️            |
| Harbor                                                      | Registry | ✔️            |
| Docker                                                      | Runtime  | ✔️            |
| Containerd                                                  | Runtime  | ✔️            |
| kubernetes                                                  | Cluster  | ✔️            |

## 🛴 Architecture
![](docs/architecture.png)

## 🏘️ Contact Us

1. You can make bug feedback and feature suggestions directly through GitHub Issues.
2. By scanning the QR code below (use wechat), you can join the discussion group of veinmind users for detailed
   discussions by adding the veinmind assistant.

![](docs/veinmind-group-qrcode.jpg)

## ✨ 404 starlink project

<img src="https://github.com/knownsec/404StarLink-Project/raw/master/logo.png" width="30%">

veinmind-tools now joined 404 starlink project (https://github.com/knownsec/404StarLink)
