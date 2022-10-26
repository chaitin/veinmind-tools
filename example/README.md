# veinmind-example

快速的初始化一个veinmind插件

+ step.1 确保当前目录处于`./veinmind-tools`的项目目录下  
+ step.2 运行 `bash ./scripts/cli/add_a_new_plugin.sh`  
+ step.3 在`plugins/go`或`plugins/python`目录下将出现对应初始化的插件目录   
+ step.4 在`cmd/cli.go`或`scan.py`函数内编写你自己的扫描逻辑  


![demo](../docs/veinmind-example/exampledemo.png)