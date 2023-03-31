<h1 align="center"> veinmind-runner </h1>

<p align="center">
veinmind-runner is a safety tool platform for pulse container developed by Changting Technology
</p>

## 📸 Basic Introduction

With rich research and development experience as the background, Changting team designed a set of plug-in system in [veinmind-sdk](). With the support of the plugin system, you only need to call the API provided by [veinmind-sdk]() to automatically generate plug-ins that meet the standard specifications. (
See [example](./example)) for a code example.
As a plugin platform, 'veinmind-runner' will automatically scan for compliant plugins and pass the image information to the corresponding plugin.
! [](https://dinfinite.oss-cn-beijing.aliyuncs.com/image/20220321150601.png)

## 🔥 Features
<b>2023-3-24 - NEW</b>
- 🔥🔥🔥 Support access to 'openai' to conduct humanized analysis of scanned security events, allowing you to more clearly understand what risks have been identified during this scan and how to operate them

> Note: When using openAI, please ensure that the current network can access openAI
> When starting a parallel container, you need to manually use docker run -e http_proxy=xxxx -e https_proxy=xxxx Set proxy (in non global proxy scenarios)


<b>Basic characteristics</b>
- Automatically scan and register plugins in the current directory (including subdirectories)
- Run different language plugins in one way
- Plugins can communicate with the 'runner' to alert on events, etc

## 💻 Compatibility

- linux/amd64
- linux/386
- linux/arm64
- linux/arm

## 🕹 Usage

### Based on executable files

Please install ` libveinmind `, installation method can refer to [official documentation] (https://github.com/chaitin/libveinmind)
#### Makefile one-click command

```
make run ARG="scan xxx"
```
#### Compile your own executable file for scanning

Compile the executable
```
make build
```
Run the executable file for scanning
```
chmod +x veinmind-runner && ./veinmind-runner scan xxx
```
### Based on the parallel container pattern
Make sure you have 'docker' and 'docker-compose' installed on your machine
#### Makefile one-click command
```
make run.docker ARG="scan xxxx"
```
#### Build your own image for scanning
Build the 'veinmind-runner' image
```
make build.docker
```
Run the container to scan
```
docker run --rm -it --mount 'type=bind,source=/,target=/host,readonly,bind-propagation=rslave' veinmind-runner scan xxx
```
### With kubernetes helm
Install 'veinmind-runner' with 'Helm' on Kubernetes to schedule scanning tasks

Please install ` Helm `, installation method can refer to [official documentation] (https://helm.sh/zh/docs/intro/install/)

Install 'veinmind-runner'
Before, can be configured to perform parameter, refer to [documents] (https://github.com/chaitin/veinmind-tools/blob/master/veinmind-runner/script/helm_chart/README.md)

Install 'veinmind-runner' using Helm

```
cd ./veinmind-runner/script/helm_chart/veinmind
helm install veinmind .
```
## ⚙ Use parameters
### Basic parameters
Refer to [veinmind-runner usage parameters documentation](docs/veinmind-runner.md)
### Advanced parameters
1.Support intelligent analysis of results using openai
> Precondition 1: Need to prepare openai_ Key, please refer to: https://platform.openai.com/account/api-keys
> Precondition 2: The network can access openai during scanning

Use the `-- analyze` parameter to add the scan to openai result analysis:

`./veinmind-runner scan image --enable-analyze --openai-token <your_openai_key>`

If you feel that the analysis result is not satisfactory, you can customize the query result statement to adjust openai's analysis of the result:

`./veinmind-runner scan image --enable-analyze --openai-token <your_openai_key> -p "Please analyze the following security incidents"`

Or：
`./veinmind-runner scan image --enable-analyze --openai-token <your_openai_key> -p "Parse what happened to the following json"`

You can also analyze the resulting file after scanning:

`./veinmind-runner analyze -r <path_to_result.json> -t <your_openai_key>`

This method will parse the `result. json` and also support customized queries with the `-p` parameter.

2. Support docker image blocking

```bash
# first
./veinmind-runner authz -c config.toml
# second
dockerd --authorization-plugin=veinmind-broker
```

The 'config.toml' contains the following fields

|        | **field name**    | **field properties** | **meanings**             |
|--------|-------------------|----------------------|--------------------------|
| policy | action            | string               | behavior need to monitor |
|        | enabled_plugins   | []string             | use which plugins        |
|        | plugin_params     | []string             | each plugin parameters   |
|        | risk_level_filter | []string             | risk level               |
|        | block             | bool                 | whether block            |
|        | alert             | bool                 | whether alarm            |
| log    | report_log_path   | string               | scan log                 |
|        | authz_log_path    | string               | block services log       |

- the action in principle support [DockerAPI] (https://docs.docker.com/engine/api/v1.41/#operation/) provides the operation interface
- Use the 'veinmind-weakpass' plugin to scan the 'ssh' server when 'creating a container' or 'pushing an image' and if a weak password is found, the risk level is' High '
This action is blocked and a warning is issued. Finally, the scan results are stored in 'plugin.log' and the risk results are stored in 'auth.log'.

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
2. Plugin custom parameters
```
/veinmind-runner scan image -- [plugin name]:[Run plugin function cmd].[Parameter name]=[custom value]
```
Examples:
```
./veinmind-runner scan image -- veinmind-weakpass:scan/image.serviceName=ssh
```
![](../docs/runner_1.jpg)