<h1 align="center"> veinmind-runner </h1>

<p align="center">
veinmind-runner is a veinmind container security tool platform developed by Chaitin Technology
</p>

## basic introduction

Chaitin team designed a plug-in system in [veinmind-sdk]() based on rich R&D experience. With the support of this plug-in system, you only need to call the API provided by [veinmind-sdk]() to automatically generate plug-ins that meet the standard specifications. (
Specific code examples can be found in [example](./example))
As a plug-in platform, `veinmind-runner` will automatically scan the plug-ins that meet the specifications, and pass the image information that needs to be scanned to the corresponding plug-in.
![](https://dinfinite.oss-cn-beijing.aliyuncs.com/image/20220321150601.png)

## Features

- Automatically scan and register plugins in the current directory (including subdirectories)
- Uniformly run Wenmai plug-ins based on different languages
- Plug-ins can communicate with `runner`, such as reporting events for alarms, etc.

## Compatibility

- linux/amd64
- linux/386
- linux/arm64
- linux/arm

## before the start

### Installation method one

Please install `libveinmind` first, the installation method can refer to [official document](https://github.com/chaitin/libveinmind)

Optionally compile `veinmind-runner` manually,
Or find the compiled `veinmind-runner` on the [Release](https://github.com/chaitin/veinmind-tools/releases) page to download

### Installation Method 2

Based on the parallel container mode, get the image of `veinmind-runner` and start it

```
docker run --rm -it --mount 'type=bind,source=/,target=/host,readonly,bind-propagation=rslave' \
-v `pwd`:/tool/resource -v /var/run/docker.sock:/var/run/docker.sock veinmind/veinmind-runner
```

Or use the script provided by the project to start

```
chmod +x parallel-container-run.sh && ./parallel-container-run.sh
```

### Installation method three

Based on the `Kubernetes` environment, use `Helm` to install `veinmind-runner`, and execute scanning tasks regularly

Please install `Helm` first, the installation method can refer to [official document](https://helm.sh/zh/docs/intro/install/)

Install `veinmind-runner`
Before, the execution parameters can be configured, please refer to [documentation](https://github.com/chaitin/veinmind-tools/blob/master/veinmind-runner/script/helm_chart/README.md)

Install `veinmind-runner` using `Helm`

```
cd ./veinmind-runner/script/helm_chart/veinmind
helm install veinmind .
```

## use

1. Scan the local image

```shell
./veinmind-runner scan image 
```

2. Scan all local images with specific runtime

```shell
./veinmind-runner scan image [dockerd:/containerd:]
```

For example: 
```shell
./veinmind-runner scan image dockerd:nginx
```

3. Scan the remote mirror. If the remote warehouse needs authentication, you need to use the -c parameter to specify the authentication information file in toml format (docker.io authentication is not supported yet)

```shell
./veinmind-runner scan image registry:server/image
```
For example:
```shell
#Scan the nginx:latest image of index.docker.io
./veinmind-runner scan image registry:nginx
```

```shell
#Scan the nginx:x.x image of index.docker.io
./veinmind-runner scan image registry:nginx:x.x
```

```shell
#Scan the veinmind image under your-registry-address warehouse
./veinmind-runner scan image registry:<your-registry-address>/veinmind/
```

```shell
#Scan the veinmind/veinmind-weakpass image under your-registry-address warehouse
./veinmind-runner scan image registry:<your-registry-address>/veinmind/veinmind-weakpass
```

```shell
#Scan the veinmind/veinmind-weakpass image under your-registry-address warehouse
./veinmind-runner scan image registry:<your-registry-address>/veinmind/veinmin-weakpass -c auth.toml
```

The format of `auth.toml` is as follows, `registry` represents the warehouse address, `username` represents the user name, and `password` represents the password or token

```toml
[[auths]]
registry = "index.docker.io"
username = "admin"
password = "password"
[[auths]]
registry = "registry.private.net"
username = "admin"
password = "password"
```

4. Scan local IaC files

```shell
./veinmind-runner scan iac host:path/to/iac-file
./veinmind-runner scan iac path/to/iac-file
```

5. Scan the IaC file of the remote git repository

```shell
./veinmind-runner scan iac git: http://xxxxxx.git
```
```shell
#auth
./veinmind-runner scan iac git:git@xxxxxx --sshkey=/your/ssh/key/path
./veinmind-runner scan iac git:http://{username}:password@xxxxxx.git
```
```shell
# add proxy
./veinmind-runner scan iac git:http://xxxxxx.git --proxy=http://127.0.0.1:8080
./veinmind-runner scan iac git:http://xxxxxx.git --proxy=scoks5://127.0.0.1:8080
```
```shell
# disable tls
./veinmind-runner scan iac git:http://xxxxxx.git --insecure-skip=true
```

6. Scan the remote kubernetes IaC configuration (you need to manually specify the kubeconfig file)

```
./veinmind-runner iac kubernetes:resource/name -n namespace --kubeconfig=/your/k8sConfig/path
```

7. Scan all local containers (if the container runtime type is not specified, it will try docker and containerd in sequence by default)

```
./veinmind-runner scan container [dockerd:/containerd:]
```

8. Scan the local container (if the container runtime type is not specified, it will try docker and containerd in sequence by default)

```
./veinmind-runner scan container [dockerd:/containerd:]containerID/containerRef
```
container runtime type

- dockerd
- containerd

9. Use `glob` to filter the plugins required to run

```
./veinmind-runner scan image -g "**/veinmind-malicious"
```

10. List the current plugin list

```
./veinmind-runner list plugin
```

11. Specify the container runtime path

```
./veinmind-runner scan image --docker-data-root [your_path]
```

```
./veinmind-runner scan image --containerd-root [your_path]
```

12. Support docker image blocking function

```bash
#first
./veinmind-runner authz -c config.toml
# second
dockerd --authorization-plugin=veinmind-broker
```

Where `config.toml` contains the following fields

| | **field name** | **field attribute** | **meaning** |
|----------|-------------------|----------|------- --|
| policy | action | string | Behavior to be monitored |
| | enabled_plugins | []string | Which plugins to use |
| | plugin_params | []string | Parameters for each plugin |
| | risk_level_filter | []string | risk level |
| | block | bool | whether to block |
| | alert | bool | whether to alarm |
| log | report_log_path | string | plugin scan log |
| | authz_log_path | string | Denial of service log |

- action in principle supports [DockerAPI](https://docs.docker.com/engine/api/v1.41/#operation/) provides the operation interface
- The following configuration means: when `creating a container` or `push image`, use the `veinmind-weakpass` plug-in to scan the `ssh` service, if a weak password is found, and the risk level is `High`
  prevents this operation and issues a warning. Finally, the scan results are stored in `plugin.log`, and the risk results are stored in `auth.log`.

```toml
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