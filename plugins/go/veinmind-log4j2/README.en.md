
<h1 align="center"> veinmind-log4j2 </h1>

<p align="center">
veinmind-log4j2 is mainly used to scan log4j jar files for CVE-2021-44228 vulnerability</p>

## Features

- Quickly scan containers/images for log4j2 risks
- Support detection of fat jars, jars containing jars, etc
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
make run ARG= "scan xxx"
```
#### Compile executable files for scanning.

Compile executable file.
```
make build.
```
Run the executable file to scan.
```
chmod + x veinmind-log4j2 & &. / veinmind-log4j2 scan xxx.
```
### Based on parallel container mode.
Make sure `docker` and `docker` and `dockere` are installed on the machine.
#### Makefile one-button command.
```
make run.docker ARG= "scan xxxx"
```
#### Build image for scanning.
Build an `veinmind-log4j2` image.
```
make build.docker.
```
Run the container to scan.
```
docker run-rm-it-mount 'type=bind,source=/,target=/host,readonly,bind-propagation=rslave' veinmind-log4j2 scan xxx.
```

## Using parameters.

1. Specify the image name or image ID and scan (the corresponding image needs to exist locally).

```
./veinmind-log4j2 scan image [imageID/imageName].
```
![](../../../docs/veinmind-log4j2/log4j2_scan_image_1.jpg)

2. Scan all local images.

```
./veinmind-log4j2 scan image.
```
![](../../../docs/veinmind-log4j2/log4j2_scan_image_2.jpg)

3. Specify the container name or container ID and scan.

```
./veinmind-log4j2 scan container [containerID/containerName].
```
![](../../../docs/veinmind-log4j2/log4j2_scan_container_1.jpg)
4. Scan all local containers.

```
./veinmind-log4j2 scan container.
```
![](../../../docs/veinmind-log4j2/log4j2_scan_container_2.jpg)
5. Specify output format

```
./veinmind-log4j2 scan container [containerID/containerName] -f html
#supported format： html,json,cli（default）
```
![](../../../docs/veinmind-log4j2/log4j2_format.jpg)