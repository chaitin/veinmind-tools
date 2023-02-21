<h1 align="center"> veinmind-backdoor </h1>

<p align="center">
veinmind-backdoor is a backdoor scanning tool for image developed by Chaitin Technology
</p>

## Features

- Quicklly scan backdoors in the image

    |  plugin | function  | 
    |---|---|
    |  crontab | scan  crontab config for backdoors|
    |  bashrc  | scan bash startup scripts for backdoors |
    |  sshd | scan for sshd softlink backdoors  |
    |  service | scan for malicious system services |
    |  tcpwrapper | scan for tcpwrapper backdoors |

- Supports writing backdoor detection scripts in plugin mode
- Support `containerd`/`dockerd` image backdoor scanning

## compatibility

- linux/amd64
- linux/386
- linux/arm64
- linux/arm

## Prepare

### install by package manager

1. install `libveinmind`  firstlly ，you can click here [offical document](https://github.com/chaitin/libveinmind) for more info

2. install python dependencies which `veinmind-backdoor` need，execute the command in the project directory
    ```
    cd ./plugins/python/veinmind-backdoor
    pip install -r requirements.txt
    ```

### install by parallel container
1. Install by Parallel Container，pull `veinmind-backdoor` iamge  and start
    ```
    docker run --rm -it --mount 'type=bind,source=/,target=/host,readonly,bind-propagation=rslave' registry.veinmind.tech/veinmind/veinmind-backdoor
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
    python scan.py --format [formattype] scan-images
    ```

    output type
    - stdout
    - json

5. Specify output path
    ```
    python scan.py --format json --output /tmp scan-images
    ```

## Demo
1. Scan the image which name is `service`
![](https://dinfinite.oss-cn-beijing.aliyuncs.com/image/20220329141342.png)

2. Scan all local images
![](https://dinfinite.oss-cn-beijing.aliyuncs.com/image/20220329141357.png)