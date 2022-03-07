<h1 align="center"> veinmind-sensitive </h1>

<p align="center">
veinmind-sensitive 是由长亭科技自研的一款镜像敏感信息扫描工具 
</p>

## 功能特性

- 快速扫描镜像中的敏感信息
- 支持敏感信息扫描规则自定义
- 支持`containerd`/`dockerd`镜像文件系统弱口令扫描

## 兼容性

- linux/amd64
- linux/386
- linux/arm64
- linux/arm

## 开始之前
请先安装`libveinmind`，安装方法可以参考[官方文档](https://github.com/chaitin/libveinmind)

然后安装`veinmind-sensitive`所需要的`python`依赖
```
pip install -r requirements.txt
```

## 使用

1.指定镜像名称或镜像ID并扫描 (需要本地存在对应的镜像)

```
python scan.py --name [imagename/imageid]
```

2.扫描所有本地镜像

```
python scan.py
```

3.指定镜像类型
```
python scan.py --engine [enginetype]
```

镜像类型
- dockerd
- containerd

4.指定输出类型
```
python scan.py --output [outputtype]
```

## 规则字段说明
- id: 规则标识符
- description: 规则描述
- match: 内容匹配规则，默认为正则
- filepath: 路径匹配规则，默认为正则
- env: 环境变量匹配规则，默认为正则且忽略大小写

## 演示
1.扫描指定镜像名称 `sensitive`
![](https://dinfinite.oss-cn-beijing.aliyuncs.com/image/20220215163700.png)

2.扫描所有镜像
![](https://dinfinite.oss-cn-beijing.aliyuncs.com/image/20220215164355.png)