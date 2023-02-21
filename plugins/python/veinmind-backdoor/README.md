<h1 align="center"> veinmind-backdoor </h1>

<p align="center">
veinmind-backdoor 是由长亭科技自研的一款镜像后门扫描工具 
</p>

## 功能特性

- 快速扫描镜像中的后门

|  插件 | 功能  | 
|---|---|
|  crontab | 扫描定时任务中是否包含后门  |
|  bashrc  | 扫描 bash 启动脚本是否包含后门  |
|  sshd | 扫描 sshd 软链接后门  |
|  service | 扫描恶意的系统服务 |
|  tcpwrapper | 扫描 tcpwrapper 后门 |

- 支持以插件模式编写后门检测脚本
- 支持`containerd`/`dockerd`镜像后门扫描

## 兼容性

- linux/amd64
- linux/386
- linux/arm64
- linux/arm

## 开始之前

### 安装方式一

请先安装`libveinmind`，安装方法可以参考[官方文档](https://github.com/chaitin/libveinmind)

然后安装`veinmind-backdoor`所需要的`python`依赖，在项目目录执行命令
```
cd ./plugins/python/veinmind-backdoor
pip install -r requirements.txt
```

### 安装方式二

基于平行容器的模式，获取 `veinmind-backdoor` 的镜像并启动
```
docker run --rm -it --mount 'type=bind,source=/,target=/host,readonly,bind-propagation=rslave' registry.veinmind.tech/veinmind/veinmind-backdoor
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

容器运行时类型
- dockerd
- containerd

4.指定输出类型
```
python scan.py --format [formattype] scan-images
```

输出类型
- stdout
- json

5.指定输出路径
```
python scan.py --format json --output /tmp scan-images
```

## 演示
1.扫描指定镜像名称 `service`
![](https://dinfinite.oss-cn-beijing.aliyuncs.com/image/20220329141342.png)

2.扫描所有镜像
![](https://dinfinite.oss-cn-beijing.aliyuncs.com/image/20220329141357.png)