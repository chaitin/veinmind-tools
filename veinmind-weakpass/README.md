<h1 align="center"> veinmind-weakpass </h1>

<p align="center">
veinmind-weakpass 是由长亭科技自研的一款镜像弱口令扫描工具 
</p>

## 功能特性

- 快速扫描镜像中的弱口令
- 支持弱口令宏定义
- 支持并发扫描弱口令
- 支持自定义用户名以及字典
- 支持`containerd`/`dockerd`镜像文件系统弱口令扫描

## 兼容性

- linux/amd64
- linux/386
- linux/arm64
- linux/arm

## 开始之前
请先安装`libveinmind`，安装方法可以参考[官方文档](https://github.com/chaitin/libveinmind)

## 使用

1.指定镜像名称或镜像ID并扫描 (需要本地存在对应的镜像)

```
./veinmind-weakpass scan [imagename/imageid]
```

2.扫描所有本地镜像

```
./veinmind-weakpass scan
```

3.指定镜像类型
```
./veinmind-weakpass scan -e [enginetype]
```

镜像类型
- dockerd
- containerd

4.指定扫描用户名类型
```
./veinmind-weakpass scan -u username
```

5.指定自定义扫描字典
```
./veinmind-weakpass scan -d ./pass.dict
```

6.解压默认字典到本地磁盘
```
./veinmind-weakpass extract
```

## 演示
1.扫描指定镜像名称 `weakpass`
![](https://dinfinite.oss-cn-beijing.aliyuncs.com/image/20220127151043.png)

2.扫描所有镜像
![](https://dinfinite.oss-cn-beijing.aliyuncs.com/image/20220127151350.png)