<h1 align="center"> veinmind-weakpass </h1>

<p align="center">
veinmind-weakpass is a weak password scanning tool for image developed by Chaitin Technology
</p>

## Features

- Quickly scan the weak password in image
- Support weak password macro definition
- Support concurrent scanning for weak passwords
- Support custom username and dictionary
- Support container runtime `containerd` and `dockerd`

## compatibility

- linux/amd64
- linux/386
- linux/arm64
- linux/arm

## Prepare

### install by package manager 

-  install `libveinmind`  firstlly ，you can click here [offical document](https://github.com/chaitin/libveinmind) for more info

### install by parallel container
- Install by Parallel Container，pull `veinmind-weakpass` iamge  and start
    ```
    docker run --rm -it --mount 'type=bind,source=/,target=/host,readonly,bind-propagation=rslave' veinmind/veinmind-weakpass scan
    ```
- or start with the script which we provided
    ```
    chmod +x parallel-container-run.sh && ./parallel-container-run.sh scan
    ```

## How to use

1. Scan image with specified image name or ID(need to have a corresponding image locally)
    ```
    ./veinmind-weakpass scan [imagename/imageid]
    ```

2. Scan all local images

    ```
    ./veinmind-weakpass scan
    ```

3. Specify container runtime type
    ```
    ./veinmind-weakpass scan --containerd
    ```

    container runtime type
    - dockerd
    - containerd

4. Specify the username which you want to scan
    ```
    ./veinmind-weakpass scan -u username
    ```

5. Specify the custom dict
    ```
    ./veinmind-weakpass scan -d ./pass.dict
    ```
6. Specify the services name
    ```
    ./veinmind-weakpass scan -a ssh,mysql,redis
    ```
    - support these service currently

        | serverName | version |
        |:----------:|:-------:|
        |     ssh    |   all   |
        |    mysql   |   8.X   |
        |    redis   |   all   |
        |   tomcat   |   all   |

7. Extract default dictionary to local disk
    ```
    ./veinmind-weakpass extract
    ```

## Demo
1.  Scan the image which name is `test` and all service supported
![](../../../docs/veinmind-weakpass/weakpasscandemo1.png)
2. Specify the image `test` and scan `ssh` service in the image
![](../../../docs/veinmind-weakpass/weakpasscandemo2.png)
2. Scan `ssh` service in all images
![](../../../docs/veinmind-weakpass/weakpasscandemo3.png)
