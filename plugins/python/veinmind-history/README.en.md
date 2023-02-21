<h1 align="center"> veinmind-history </h1>

<p align="center">
veinmind-history is an image anomaly history command scanning tool developed by Chaitin Technology
</p>

## Features

- Quickly scan the image for abnormal history commands
- Support custom historical command detection rules
- Support two container runtime `containerd` and `dockerd`

## Compatibility

- linux/amd64
- linux/386
- linux/arm64
- linux/arm

## Prepare

### install by package manager

1. install `libveinmind`  firstlly ，you can click here [offical document](https://github.com/chaitin/libveinmind) for more info

2. install python dependencies which `veinmind-history` need，execute the command in the project directory
    ```
    cd ./plugins/python/veinmind-history
    pip install -r requirements.txt
    ```

### install by parallel container
1. Install by Parallel Container，pull `veinmind-history` iamge  and start
    ```
    docker run --rm -it --mount 'type=bind,source=/,target=/host,readonly,bind-propagation=rslave' registry.veinmind.tech/veinmind/veinmind-history
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
- description: description for rules
- instruct: click [here](https://docs.docker.com/engine/reference/builder/) for more info 
- match: content matching rules, the default is regular

## Demo
1. Scan the image which name is `history`
![](https://dinfinite.oss-cn-beijing.aliyuncs.com/image/20220329111927.png)
2. Scan all images
![](https://dinfinite.oss-cn-beijing.aliyuncs.com/image/20220329111948.png)