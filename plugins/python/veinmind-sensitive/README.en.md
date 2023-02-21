<h1 align="center"> veinmind-sensitive </h1>

<p align="center">
veinmind-sensitive is an image sensitive information scanning tool developed by Chaitin Technology 
</p>

## Features

- Quickly scan images for sensitive information
- Support custom sensitive information scanning rules
- Support container runtime `containerd` and `dockerd`

## Compatibility

- linux/amd64
- linux/386
- linux/arm64
- linux/arm

## Prepare

### install by package manager

1. install `libveinmind`  firstlly ，you can click here [offical document](https://github.com/chaitin/libveinmind) for more info

2. install python dependencies which `veinmind-sensitive` need，execute the command in the project directory

   ```
   cd ./plugins/python/veinmind-sensitive
   pip install -r requirements.txt
   ```
### install by parallel container

1. Install by Parallel Container，pull `veinmind-sensitive` iamge  and start
    ```
    docker run --rm -it --mount 'type=bind,source=/,target=/host,readonly,bind-propagation=rslave' registry.veinmind.tech/veinmind/veinmind-sensitive
    ```
2. or start with the script which we provided
    ```
    chmod +x parallel-container-run.sh && ./parallel-container-run.sh
    ```

## How to use

1. Scan image with specified image name or ID(need to have a corresponding image locally)

    ```
    python scan.py scan-images [imagename/imageid]
    ```

2. Scan all local images

    ```
    python scan.py scan-images
    ```

3. Specify the container runtime type
    ```
    python scan.py scan-images --containerd
    ```

    container runtime type
    - dockerd
    - containerd

4. Specify output type
    ```
    python scan.py --output [outputtype] scan-images
    ```

## Rule Field Description
- id: rule identifier
- description: rule description
- match: content matching rules, the default is regular
- filepath: path matching rules, the default is regular
- env: environment variable matching rules, the default is regular and ignores case

## Demo
1. Scan the image which name is `sensitive`
![](https://dinfinite.oss-cn-beijing.aliyuncs.com/image/20220329142155.png)

2. Scan all local images
![](https://dinfinite.oss-cn-beijing.aliyuncs.com/image/20220329142506.png)