<h1 align="center"> veinmind-host </h1>

<p align="center">
veinmind-host 是由长亭科技自研的一款用于运行和管理问脉插件的宿主
</p>

## 功能特性

- 自动扫描并注册当前目录下(含子目录)的插件
- 统一运行基于不同语言实现的问脉插件
- 支持 `dockerd`/`containerd` 两种容器运行时

## 兼容性

- linux/amd64
- linux/386
- linux/arm64
- linux/arm

## 使用

1.指定镜像名称或镜像ID并扫描 (需要本地存在对应的镜像)

```
./veinmind-host scan [imagename/imageid]
```

2.扫描所有本地镜像

```
./veinmind-host scan
```

3.指定容器运行时类型
```
./veinmind-host scan --containerd
```

容器运行时类型
- dockerd
- containerd

4.使用`glob`筛选插件
```
./veinmind-host scan -g "**/veinmind-malicious"
```

## 演示
1.扫描指定镜像名称 `xmrig/xmrig`
![](https://dinfinite.oss-cn-beijing.aliyuncs.com/image/20220314150819.png)