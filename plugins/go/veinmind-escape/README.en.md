
<h1 align="center"> veinmind-escape </h1>

<p align="center">
veinmind-escape is an escape risk scanning tool developed by Changting Technology.
</p>

## Features

- quickly scan containers / images for escape risks.
- supports `docker` / `containerd` container runtime.
- multiple report formats such as `JSON` / `CLI` / `HTML` are supported

## Compatibility

- linux/amd64
- linux/386
- linux/arm64
- linux/arm

## # How to use.

### Based on executable file.

Please install `libveinmind` first. For installation method, please see [official documentation] (https://github.com/chaitin/libveinmind)).
#### Makefile one-button command.

```
Make run ARG= "scan xxx"
```
#### Compile executable files for scanning.

Compile executable file.
```
Make build.
```
Run the executable file to scan.
```
Chmod + x veinmind-escape & &. / veinmind-escape scan xxx.
```
### Based on parallel container mode.
Make sure `docker` and `docker` and `dockere` are installed on the machine.
#### Makefile one-button command.
```
Make run.docker ARG= "scan xxxx"
```
#### Build image for scanning.
Build an `veinmind- escalate` image.
```
Make build.docker.
```
Run the container to scan.
```
Docker run-rm-it-mount 'type=bind,source=/,target=/host,readonly,bind-propagation=rslave' veinmind-escape scan xxx.
```

## Using parameters.

1. Specify the image name or image ID and scan (the corresponding image needs to exist locally).

```
. / veinmind-escape scan image [imageID/imageName].
```
![](../../../docs/veinmind-escape/veinmind-escape_scan_image_01.jpg)

2. Scan all local images.

```
. / veinmind-escape scan image.
```
![](../../../docs/veinmind-escape/veinmind-escape_scan_image_02.jpg)

3. Specify the container name or container ID and scan.

```
. / veinmind-escape scan container [containerID/containerName].
```
![](../../../docs/veinmind-escape/veinmind-escape_scan_container_01.jpg)


4. Scan all local containers.

```
. / veinmind-escape scan container.
```
![](../../../docs/veinmind-escape/veinmind-escape_scan_container_02)

5. Specify output format

```
./veinmind-escape scan container [containerID/containerName] -f html
#supported format： html,json,cli（default）
```
![](../../../docs/veinmind-escape/veinmind-escape_format.jpg)