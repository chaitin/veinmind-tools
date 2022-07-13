# Helm chart for Kubernetes

veinmind-runner镜像定时启动脚本，可使用crontab语法配置定时执行扫描镜像

## 安装
1. 确定本地运行Kubernetes服务
```bash
[root@localhost veinmind]# kubectl get pods -n kube-system
NAME                                       READY   STATUS    RESTARTS   AGE
calico-kube-controllers-6d75fbc96d-2d67s   1/1     Running   0          48m
calico-node-47fzd                          1/1     Running   0          48m
calico-typha-6576ff658-xsbbv               1/1     Running   0          48m
......
```
2. 安装helm
```
wget https://get.helm.sh/helm-v3.9.0-linux-amd64.tar.gz
tar -zxvf helm-v3.9.0-linux-amd64.tar.gz
mv linux-amd64/helm /usr/local/bin/helm
```

```bash
[root@localhost veinmind]# helm
The Kubernetes package manager

Common actions for Helm:

- helm search:    search for charts
- helm pull:      download a chart to your local directory to view
- helm install:   upload the chart to Kubernetes
- helm list:      list releases of charts
```

3. 进入`helm_chart\veinmind\`:
```bash
# 安装
[root@localhost veinmind]# helm install veinmind .
# 卸载
[root@localhost veinmind]# helm uninstall veinmind
```

## 配置解析
项目主要配置信息位于`values.ymal`:
```ymal
jobs:
  ### REQUIRED ###
  - name: veinmind-runner
    image:
      repository: veinmind/veinmind-runner
      tag: latest
      imagePullPolicy: IfNotPresent
    schedule: "0 */8 * * *"   ### 扫描周期配置
    failedJobsHistoryLimit: 1
    successfulJobsHistoryLimit: 3
    concurrencyPolicy: Allow
    restartPolicy: OnFailure
  ### OPTIONAL ###
    command: ["/tool/entrypoint.sh"] ### 程序入口点
    args:
      - "scan-host"     ### 运行参数
    resources:          ### 资源配置,1000m == 1 个 CPU 单元，相当于1 个物理 CPU 核，或1 个虚拟核
      limits:
        cpu: 1000m
        memory: 256Mi
      requests:
        cpu: 1000m
        memory: 256Mi
    volumes:
      - name: files-mount
        hostPath:
          path: /
      - name: sock-path
        hostPath:
          path: /var/run/docker.sock
    volumeMounts:
      - name: files-mount
        mountPath: /host
      - name: sock-path
        mountPath: /var/run/docker.sock
```

## 运行截图
![img.png](img/KuboardView.png)

![img.png](img/kubctl.png)

扫描结果请查询日志
![img.png](img/logs.png)