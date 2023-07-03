<h1 align="center"> veinmind-trace </h1>

<p align="center">
veinmind-trace 是由长亭科技自研的一款容器安全检测工具
</p>

## 功能特性
+ 快速扫描容器中的异常进程:
  1. 隐藏进程(mount -o bind方式)
  2. 反弹shell的进程
  3. 带有挖矿、黑客工具、可疑进程名的进程
  4. 包含 Ptrace 的进程
+ 快速扫描容器中的异常文件系统: 
  1. 敏感目录权限异常
  2. cdk 工具利用痕迹检测
+ 快速扫描容器中的异常用户: 
  1. uid=0 的非root账户
  2. gid=0 的非root账户
  3. uid相同的用户
+ 支持`containerd`/`dockerd`容器运行时