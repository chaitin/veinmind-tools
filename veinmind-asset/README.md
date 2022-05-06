<h1 align="center"> veinmind-asset </h1>

<p align="center">
veinmind-asset 主要用于扫描容器镜像的内资产信息
</p>

## 功能特性

- 扫描镜像的OS信息
- 扫描镜像内系统安装的packages
- 扫描镜像内应用安装的libraries

## 使用

1.指定镜像名称或镜像ID并扫描 (需要本地存在对应的镜像)

```
./veinmind-asset scan [imagename/imageid]
```
![](https://cdn.dvkunion.cn/16510316433810.jpg)

2.扫描本地全部镜像

```
./veinmind-asset scan
```

3.输出详细结果
```
./veinmind-asset scan -v
```
![](https://cdn.dvkunion.cn/16510317401391.jpg)

4.输出指定类型的详细结果
```
./veinmind-asset scan -v --type [os/python/jar/pip/npm.......]
```
![](https://cdn.dvkunion.cn/16510559474726.jpg)

5.输出详细结果到文件

```
./veinmind-asset scan -f [csv/json]
```
![](https://cdn.dvkunion.cn/16510318063574.jpg)