package filter

type ScriptSuffix string

const (
	PHP_SUFFIX ScriptSuffix = ".php"
	JSP_SUFFIX ScriptSuffix = ".jsp"
	ASP_SUFFIX ScriptSuffix = ".asp"
)

type ScriptType string

const (
	UNKNOWN_TYPE            = "UNKNOWN"
	PHP_TYPE     ScriptType = "php"
	JSP_TYPE     ScriptType = "jsp"
	ASP_TYPE     ScriptType = "asp"
)

var scriptSuffixes []ScriptSuffix = func() []ScriptSuffix {
	return []ScriptSuffix{
		PHP_SUFFIX,
		JSP_SUFFIX,
		ASP_SUFFIX,
	}
}()

var scriptSuffixTypeMap map[ScriptSuffix]ScriptType = func() map[ScriptSuffix]ScriptType {
	return map[ScriptSuffix]ScriptType{
		PHP_SUFFIX: PHP_TYPE,
		JSP_SUFFIX: JSP_TYPE,
		ASP_SUFFIX: ASP_TYPE,
	}
}()
