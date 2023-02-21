<h1 align="center"> veinmind-history </h1>

<p align="center">
veinmind-history 是由长亭科技自研的一款镜像异常历史命令扫描工具 
</p>

## 功能特性

- 快速扫描镜像中的异常历史命令
- 支持自定义历史命令检测规则
- 支持`containerd`/`dockerd`两种容器运行时

## 兼容性

- linux/amd64
- linux/386
- linux/arm64
- linux/arm

## 开始之前

### 安装方式一

请先安装`libveinmind`，安装方法可以参考[官方文档](https://github.com/chaitin/libveinmind)

然后安装`veinmind-history`所需要的`python`依赖，在项目目录执行命令
```
cd ./plugins/python/veinmind-history
pip install -r requirements.txt
```

### 安装方式二

基于平行容器的模式，获取 `veinmind-history` 的镜像并启动
```
docker run --rm -it --mount 'type=bind,source=/,target=/host,readonly,bind-propagation=rslave' registry.veinmind.tech/veinmind/veinmind-history
```

或者使用项目提供的脚本启动
```
chmod +x parallel-container-run.sh && ./parallel-container-run.sh
```

## 使用

1.指定镜像名称或镜像ID并扫描 (需要本地存在对应的镜像)

```
python scan.py scan-images [imagename/imageid]
```

2.扫描所有本地镜像

```
python scan.py scan-images
```

3.指定容器运行时类型
```
python scan.py scan-images --containerd
```

镜像类型
- dockerd
- containerd

4.指定输出类型
```
python scan.py --output [outputtype] scan-images
```

## 规则字段说明
- description: 规则描述
- instruct: [操作指令](https://docs.docker.com/engine/reference/builder/)
- match: 内容匹配规则，默认为正则

## 演示
1.扫描指定镜像名称 `history`
![](https://dinfinite.oss-cn-beijing.aliyuncs.com/image/20220329111927.png)
2.扫描所有镜像
![](https://dinfinite.oss-cn-beijing.aliyuncs.com/image/20220329111948.png)