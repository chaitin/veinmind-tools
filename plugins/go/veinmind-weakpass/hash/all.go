package hash

// 注册所有的hash方法
// service模块可以通过hash的ID获取具体的hash实例
var All = []Hash{
	&Plain{},
	&MysqlNative{},
	&Shadow{},
}
