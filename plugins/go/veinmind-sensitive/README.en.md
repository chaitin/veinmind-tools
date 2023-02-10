<h1 align="center"> veinmind-sensitive </h1>

<p align="center">
veinmind-sensitive is a mirror sensitive information scanning tool developed by Changting Technology
</p>

## Features

- Quickly scan the image for sensitive information
- Support sensitive information scanning rule customization
- Supports weak password scanning of 'containerd'/' dockerd 'image file systems

## Compatibility

- linux/amd64
- linux/386
- linux/arm64
- linux/arm

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
chmod +x veinmind-sensitive && ./veinmind-sensitive scan xxx
```
### Based on the parallel container pattern
Make sure you have 'docker' and 'docker-compose' installed on your machine
#### Makefile one-click command
```
make run.docker ARG="scan xxxx"
```
#### Build your own image for scanning
Build the 'veinmind-sensitive' image
```
make build.docker
```
Run the container to scan
```
docker run --rm -it --mount 'type=bind,source=/,target=/host,readonly,bind-propagation=rslave' veinmind-sensitive scan  xxx
```

## Use parameters

1. Specify the image name or image ID and scan (if the image exists locally)

```
./veinmind-sensitive scan image [imagename/imageid]
```
![](../../../docs/veinmind-sensitive/sensitive-01.jpeg)

2. Scan all local images

```
./veinmind-sensitive scan image
```
![](../../../docs/veinmind-sensitive/sensitive-02-1.jpg)

![](../../../docs/veinmind-sensitive/sensitive-02-2.jpg)

Specify the output type
Supported output formats:
- html
- json
- cli (default)
```
./veinmind-sensitive scan image [imageID/imageName] -f html
```
The resulting result.html looks like this:
![](../../../docs/veinmind-sensitive/sensitive-03.jpg)

## Rule Field Description

-id: Rule identifier
-description: Description of the rule
-match: This is a content-matching rule, which defaults to regular
-filepath: Path matching rule (regular by default
-env: Environment variable matching rules; defaults to regular and ignores case