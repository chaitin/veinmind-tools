package hash

var All = []Hash{
	&Plain{},
	&MysqlNative{},
	&Shadow{},
}
