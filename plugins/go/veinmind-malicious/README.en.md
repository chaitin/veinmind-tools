<h1 align="center"> veinmind-malicious </h1>

<p align="center">
veinmind-malicious is a malicious file scanning tool for images developed by Chaitin Technology 
</p>

## Features

- Quickly scan images for malicious files(`ClamAV` and `VirusTotal` have been supported )
- Support container runtime such as `docker` / `containerd` 
- Support different output type like `JSON` / `CSV` / `HTML`

## Compatibility

- linux/amd64
- linux/386
- linux/arm64
- linux/arm

## Prepare

### install by package manager 

1. install `libveinmind`  firstlly ，you can click here [offical document](https://github.com/chaitin/libveinmind) for more info

2. make sure `docker` and `docker-compose` were installed and then start `ClamAV`。

    ```
    chmod +x veinmind-malicious && ./veinmind-malicious extract && cd scripts && docker-compose pull && docker-compose up -d
    ```

3. if you use `VirusTotal`，you should modify `scripts/.env` in which add `VT_API_KEY`
    ```
    export VT_API_KEY=xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx
    ```

### install by parallel container

1. Install by Parallel Container，pull `veinmind-malicious` iamge  and start
    ```
    docker run --rm -it --mount 'type=bind,source=/,target=/host,readonly,bind-propagation=rslave' -v `pwd`:/tool/data veinmind/veinmind-malicious scan
    ```

2. or start with the script which we provided
    ```
    chmod +x parallel-container-run.sh && ./parallel-container-run.sh scan
    ```

## How to use

1. Scan image with specified image name or ID(need to have a corresponding image locally)

    ```
    ./veinmind-malicious scan [imagename/imageid]
    ```

2. Scan all local images

```
./veinmind-malicious scan
```

3. Specify the output type (now support html/csv/json)

    ```
    ./veinmind-malicious scan -f [html/csv/json]
    ```

4. Specify the output file name

    ```
    ./veinmind-malicious scan -n [reportname]
    ```

5. Sepcify the output path

    ```
    ./veinmind-malicious scan -o [outputpath]
    ```

6. Specify the container runtime type
    ```
    ./veinmind-malicious scan --containerd
    ```

    container runtime types
    - dockerd
    - containerd

## Demo
1. Scan the image which name is `xmrig/xmrig`
![](https://dinfinite.oss-cn-beijing.aliyuncs.com/image/20220119111800.png)

2. Scan specified image ID `sha256:ba6acccedd2923aee4c2acc6a23780b14ed4b8a5fa4e14e252a23b846df9b6c1`
![](https://dinfinite.oss-cn-beijing.aliyuncs.com/image/20220119112217.png)

3. Specifiy the output path and output file name
![](https://dinfinite.oss-cn-beijing.aliyuncs.com/image/20220119112058.png)

## Report
![](https://dinfinite.oss-cn-beijing.aliyuncs.com/image/20220119142131.png)