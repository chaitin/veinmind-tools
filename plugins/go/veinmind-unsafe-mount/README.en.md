<h1 align="center"> veinmind-unsafe-mount </h1>

<p align="center">
veinmind-unsafe-mount is a container unsafe mount directory scanning tool developed by Changting Technology
</p>

## Features

- Quickly scan containers for unsafe mounts
- Support for the 'containerd'/' dockerd 'container runtime

## Compatibility

- linux/amd64
- linux/386
- linux/arm64

## Usage

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
chmod +x veinmind-unsafe-mount && ./veinmind-unsafe-mount scan xxx
```
### Based on the parallel container pattern
Make sure you have 'docker' and 'docker-compose' installed on your machine
#### Makefile one-click command
```
make run.docker ARG="scan xxxx"
```
#### Build your own image for scanning
Build the 'veinmind-unsafe-mount' image
```
make build.docker
```
Run the container to scan
```
docker run --rm -it --mount 'type=bind,source=/,target=/host,readonly,bind-propagation=rslave' veinmind-unsafe-mount  scan xxx
```

## Use parameters

1. Specify the container name or ID and scan (if the corresponding container exists locally)
```
./veinmind-unsafe-mount scan container [containerID/containerName]
```
![](../../../docs/veinmind-unsafe-mount/unsafemount_scan_container_01.jpg)
2. Scan all local containers
```
./veinmind-unsafe-mount scan container
```
![](../../../docs/veinmind-unsafe-mount/unsafemount_scan_container_02.jpg)
3. Specify the output format
Supported output formats:
- html
- json
- cli (default)
```
./veinmind-unsafe-mount scan container [containerID/containerName] -f html
```
The resulting result.html looks like this:
![](../../../docs/veinmind-unsafe-mount/unsafemount_scan_container_03.jpg)