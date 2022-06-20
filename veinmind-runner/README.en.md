<h1 align="center"> veinmind-runner </h1>

<p align="center">
veinmind-runner it's a container security host developed by Chaitin Technology.
</p>

## Introduce
With the background of rich R&D experience, the chaitin team designed a plug-in system in veinmind-sdk.
With the support of this plugin system, you only need to call the API provided by veinmind-sdk to automatically generate plugins that conform to standard specifications. (For specific code examples, see [example](./example))
As a plugin platform, `veinmind-runner` will automatically scan the plugins that conform to the specification, and pass the image information that needs to be scanned to the corresponding plugins.
![](https://dinfinite.oss-cn-beijing.aliyuncs.com/image/20220321150601.png)

## Feature

- Automatically scan and register plugins in the current directory (including subdirectories)
- Unified operation of plug-ins implemented in different languages
- Plugins can communicate with `runner`, such as reporting events for alarming, etc.

## Compatibility

- linux/amd64
- linux/386
- linux/arm64
- linux/arm

## Install

### Install by package manager

please install `libveinmind`，here is [official document](https://github.com/chaitin/libveinmind)

you can compile manually `veinmind-runner`，
or download from [Release](https://github.com/chaitin/veinmind-tools/releases)

### Install by parallel container

based on the parallel container mode, get the image of `veinmind-runner` and start it
```
docker run --rm -it --mount 'type=bind,source=/,target=/host,readonly,bind-propagation=rslave' \
-v `pwd`:/tool/resource -v /var/run/docker.sock:/var/run/docker.sock veinmind/veinmind-runner
```

or use the script provided by the project to start
```
chmod +x parallel-container-run.sh && ./parallel-container-run.sh
```

## Usage

1.specify the image name or image ID and scan (need to have a corresponding image locally)

```
./veinmind-runner scan-host [imagename/imageid]
```

2.scan all local images

```
./veinmind-runner scan-host
```

3.scan the `centos` image in the remote repository (the default is `index.docker.io` if the repository is not specified)

```
./veinmind-runner scan-registry centos
```

4.scan `registry.private.net/library/nginx` image in the remote private registry, where `auth.toml` is the authentication information configuration file, which contains the corresponding authentication information

```
./veinmind-runner scan-registry -c auth.toml registry.private.net/library/nginx
```

the format of `auth.toml` is as follows, `registry` represents the repository address, `username` represents the username, `password` represents the password or token
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

5.specify the container runtime type

```
./veinmind-runner scan-host --containerd
```

container runtime type
- dockerd
- containerd

6.filtering with `glob` requires running the plugin
```
./veinmind-runner scan-host -g "**/veinmind-malicious"
```

7.list plugin
```
./veinmind-runner list plugin
```

8.specify the container runtime path
```
./veinmind-runner scan-host --docker-data-root [your_path]
```
```
./veinmind-runner scan-host --containerd-root [your_path]
```