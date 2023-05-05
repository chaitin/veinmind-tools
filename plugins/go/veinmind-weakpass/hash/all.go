package hash

// All 注册所有的hash方法 service模块可以通过hash的ID获取具体的hash实例
var All = []Hash{
	&Plain{},
	&MySQL{},
	&MysqlNative{},         // This is the one we want to remove
	&CachingSha2Password{}, // This is the one we want to remove
	&Shadow{},
}
