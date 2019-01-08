package smart

type Expression interface {
	// 执行表达式
	Eval(expr string, params map[string]interface{}) bool
}