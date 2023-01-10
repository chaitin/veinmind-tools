<h1 align="center"> veinmind-runner </h1>

<p align="center">
veinmind-runner 是由长亭科技自研的一款问脉容器安全工具平台
</p>

## 基本介绍

长亭团队以丰富的研发经验为背景， 在 [veinmind-sdk]() 中设计了一套插件系统。 在该插件系统的支持下，只需要调用 [veinmind-sdk]() 所提供的API，即可自动化的生成符合标准规范的插件。(
具体代码示例可查看[example](./example))
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

## 开始之前

### 安装方式一

请先安装`libveinmind`，安装方法可以参考[官方文档](https://github.com/chaitin/libveinmind)

可以选择手动编译 `veinmind-runner`，
或者在[Release](https://github.com/chaitin/veinmind-tools/releases)页面中找到已经编译好的 `veinmind-runner` 进行下载

### 安装方式二

基于平行容器的模式，获取 `veinmind-runner` 的镜像并启动

```
docker run --rm -it --mount 'type=bind,source=/,target=/host,readonly,bind-propagation=rslave' \
-v `pwd`:/tool/resource -v /var/run/docker.sock:/var/run/docker.sock veinmind/veinmind-runner
```

或者使用项目提供的脚本启动

```
chmod +x parallel-container-run.sh && ./parallel-container-run.sh
```

### 安装方式三

基于`Kubernetes`环境，使用`Helm`安装`veinmind-runner`，定时执行扫描任务

请先安装`Helm`， 安装方法可以参考[官方文档](https://helm.sh/zh/docs/intro/install/)

安装`veinmind-runner`
之前，可配置执行参数，可参考[文档](https://github.com/chaitin/veinmind-tools/blob/master/veinmind-runner/script/helm_chart/README.md)

使用`Helm`安装 `veinmind-runner`

```
cd ./veinmind-runner/script/helm_chart/veinmind
helm install veinmind .
```

## 使用

1.扫描本地镜像

```
./veinmind-runner image [docker/containerd]:imageID/imageRef
```

2.扫描所有本地镜像

```
./veinmind-runner image [docker/containerd]:
```

3.扫描远程镜像，若远程仓库需要认证需使用-c参数指定toml格式的认证信息文件（暂不支持docker.io的认证）

```
./veinmind-runner image registry:server/imageRef 如果不指定server，则默认为docker.io
如：
./veinmind-runner image registry:(docker.io/)nginx                                 扫描docker.io的nginx镜像
./veinmind-runner image registry:(docker.io/)bitnami/nginx                         扫描docker.io的bitnami/nginx镜像
./veinmind-runner image registry:(docker.io/)bitnami/                              扫描docker.io中bitnami/下的所有镜像
./veinmind-runner image registry:<your-registry-address>/veinmind-weakpass         扫描your-registry-address仓库下的veinmind-weakpass镜像
./veinmind-runner image registry:<your-registry-address>/veinmind/veinmin-weakpass 扫描your-registry-address仓库下的veinmind/veinmind-weakpass镜像
./veinmind-runner image registry:<your-registry-address>/veinmind/veinmin-weakpass 扫描your-registry-address仓库下的veinmind/veinmind-weakpass镜像
./veinmind-runner image registry:<your-registry-address>/veinmind/                 扫描your-registry-address仓库中veinmind/下的镜像
```
`auth.toml` 的格式如下， `registry` 代表仓库地址， `username` 代表用户名， `password` 代表密码或 token

```
[[auths]]
	registry = "index.docker.io"
	username = "admin"
	password = "password"
[[auths]]
	registry = "registry.private.net"
	username = "admin"
	password = "password"
```

4.扫描本地IaC文件

```
./veinmind-runner iac host:path/to/iac-file
./veinmind-runner iac path/to/iac-file
```

5.扫描远端 git 仓库的 IaC 文件

```
./veinmind-runner iac git:http://xxxxxx.git 
# auth
./veinmind-runner iac git:git@xxxxxx --sshkey=/your/ssh/key/path
./veinmind-runner iac git:http://{username}:password@xxxxxx.git
# add proxy
./veinmind-runner iac git:http://xxxxxx.git --proxy=http://127.0.0.1:8080
./veinmind-runner iac git:http://xxxxxx.git --proxy=scoks5://127.0.0.1:8080
# disable tls
./veinmind-runner iac git:http://xxxxxx.git --insecure-skip=true
```

6.扫描远端 kubernetes IaC 配置(需要手动指定kubeconfig file)

```
./veinmind-runner iac kubernetes:resource/name -n namespace --kubeconfig=/your/k8sConfig/path
```

7.扫描本地所有容器

```
./veinmind-runner containerd [docker/containerd]:
```

8.扫描本地容器(容器运行时类型未指定的情况下默认为docker)

```
./veinmind-runner container [docker/containerd]:containerID/containerRef
```
容器运行时类型

- dockerd
- containerd

9.使用`glob`筛选需要运行插件

```
./veinmind-runner  image -g "**/veinmind-malicious"
```

10.列出当前插件列表

```
./veinmind-runner list plugin
```

11.指定容器运行时路径

```
./veinmind-runner image --docker-data-root [your_path]
```

```
./veinmind-runner image --containerd-root [your_path]
```

12.支持 docker 镜像阻断功能

```bash
# first
./veinmind-runner authz -c config.toml 
# second
dockerd --authorization-plugin=veinmind-broker
```

其中`config.toml`,包含如下字段

|  | **字段名**           | **字段属性** | **含义**  |
|----------|-------------------|----------|---------|
| policy   | action            | string   | 需要监控的行为 |
|          | enabled_plugins   | []string | 使用哪些插件  |
|          | plugin_params     | []string | 各个插件的参数 |
|          | risk_level_filter | []string | 风险等级    |
|          | block             | bool     | 是否阻断    |
|          | alert             | bool     | 是否报警    |
| log      | report_log_path   | string   | 插件扫描日志  |
|          | authz_log_path    | string   | 阻断服务日志  |

- action 原则上支持[DockerAPI](https://docs.docker.com/engine/api/v1.41/#operation/)所提供的操作接口
- 如下的配置表示：当 `创建容器`或`推送镜像` 时，使用 `veinmind-weakpass` 插件扫描`ssh`服务，如果发现有弱密码存在，并且风险等级为 `High`
  则阻止此操作，并发出警告。最终将扫描结果存放至`plugin.log`,将风险结果存放至`auth.log`。

``` toml
[log]
plugin_log_path = "plugin.log"
auth_log_path = "auth.log"
[listener]
listener_addr = "/run/docker/plugins/veinmind-broker.sock"
[[policies]]
action = "container_create"
enabled_plugins = ["veinmind-weakpass"]
plugin_paramas = ["veinmind-weakpass:scan.serviceName=ssh"]
risk_level_filter = ["High"]
block = true
alert = true
[[policies]]
action = "image_push"
enabled_plugins = ["veinmind-weakpass"]
plugin_params = ["veinmind-weakpass:scan.serviceName=ssh"]
risk_level_filter = ["High"]
block = true
alert = true
[[policies]]
action = "image_create"
enabled_plugins = ["veinmind-weakpass"]
plugin_params = ["veinmind-weakpass:scan.serviceName=ssh"]
risk_level_filter = ["High"]
block = true
alert = true
```