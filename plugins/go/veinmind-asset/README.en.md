<h1 align="center"> veinmind-asset </h1>

<p align="center">
veinmind-asset is mainly used to scan the internal asset information of images and containers
</p>

## Features

- Scan image OS information
- Scan the packages information installed in the image
- Scan the libraries installed by the application in the image

## How to use

1. Scan the image which name is [imagename/imageid]

    ```
    ./veinmind-asset scan [imagename/imageid]
    ```
    ![](https://cdn.dvkunion.cn/16510316433810.jpg)

2. Scan all local images

    ```
    ./veinmind-asset scan
    ```

3. Scan image with detailed results
    ```
    ./veinmind-asset scan -v
    ```
    ![](https://cdn.dvkunion.cn/16510317401391.jpg)

4. Scan image with detailed and specified type results
    ```
    ./veinmind-asset scan -v --type [os/python/jar/pip/npm.......]
    ```
    ![](https://cdn.dvkunion.cn/16510559474726.jpg)

5. Output detailed results to file

    ```
    ./veinmind-asset scan -f [csv/json]
    ```
    ![](https://cdn.dvkunion.cn/16510318063574.jpg)