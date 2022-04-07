<h1 align="center"> veinmind-runner </h1>

<p align="center">
veinmind-runner 是由长亭科技自研的一款问脉容器安全工具平台
</p>

## 基本介绍
长亭团队以丰富的研发经验为背景， 在 [veinmind-sdk]() 中设计了一套插件系统。
在该插件系统的支持下，只需要调用 [veinmind-sdk]() 所提供的API，即可自动化的生成符合标准规范的插件。(具体代码示例可查看[example](./example))
`veinmind-runner`作为插件平台，会自动化的扫描符合规范的插件，并将需要扫描的镜像信息传递给对应的插件。
![](https://dinfinite.oss-cn-beijing.aliyuncs.com/image/20220321150601.png)

## 功能特性

- 自动扫描并注册当前目录下(含子目录)的插件
- 统一运行基于不同语言实现的问脉插件
- 插件可以和`runner`进行通信，如上报事件进行告警等

## 兼容性

- linux/amd64
- linux/386
- linux/arm64
- linux/arm

## 使用

1.指定镜像名称或镜像 ID 并扫描 (需要本地存在对应的镜像)

```
./veinmind-runner scan-host [imagename/imageid]
```

2.扫描所有本地镜像

```
./veinmind-runner scan-host
```

3.扫描远程仓库中的`centos`镜像(不指定仓库默认为`index.docker.io`)

```
./veinmind-runner scan-registry -r "centos"
```

4.扫描远程私有仓库`registry.private.net`中的`nginx`镜像，其中用户名为`admin`，密码为`password`

```
./veinmind-runner scan-registry --address registry.private.net --reponames "nginx" \
--username admin  --password password
```

5.指定容器运行时类型

```
./veinmind-runner scan-host --containerd
```

容器运行时类型
- dockerd
- containerd

6.使用`glob`筛选插件
```
./veinmind-runner scan-host -g "**/veinmind-malicious"
```
