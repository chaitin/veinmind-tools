# veinmind-escalate  

<h1 align="center"> veinmind-escalate </h1>

<p align="center">
veinmind-malicious is an escape risk scanning tool developed by Changting Technology
</p>

## Features

- Quickly scan containers for escape risks
- Supports the 'docker'/' containerd 'container runtime
- Support JSON/CSV/HTML report formats

## Compatibility

- linux/amd64
- linux/386
- linux/arm64
- linux/arm

## Before we begin

### Installation 1

Please install ` libveinmind `, installation method can refer to [official documentation] (https://github.com/chaitin/libveinmind)

Make sure you have 'docker' and 'docker-compose' installed on your machine, and start 'ClamAV'.

` ` `
chmod +x veinmind-escalate && ./veinmind-escalte extract && cd scripts && docker-compose pull && docker-compose up -d
` ` `

If you're using VirusTotal, you'll need to declare the 'VT_API_KEY' in an environment variable or in your 'scripts/.env' file
` ` `
export VT_API_KEY=xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx
` ` `

### Installation 2

Take an image of 'veinmind-escalate' and launch it based on the parallel container pattern
` ` `
docker run --rm -it --mount 'type=bind,source=/,target=/host,readonly,bind-propagation=rslave' -v `pwd`:/tool/data  veinmind/veinmind-escalate scan
` ` `

Or use a script provided by the project
` ` `
chmod +x parallel-container-run.sh && ./parallel-container-run.sh scan
` ` `

## Usage

1. Specify the image name or image ID and scan (if the image exists locally)

` ` `
./veinmind-escalate scan image [imageID/imageName]
` ` `

2. Scan all local images

` ` `
./veinmind-escalate scan image
` ` `

3. Specify the container name or container ID and scan

` ` `
./veinmind-escalate scan container [imageID/imageName]
` ` `

4. Scan all local containers

` ` `
./veinmind-escalate scan container
` ` `
