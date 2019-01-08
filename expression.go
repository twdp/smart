package smart

import (
	"github.com/Knetic/govaluate"
	"github.com/astaxie/beego/logs"
)

type Expression interface {

	// 执行表达式
	Eval(expr string, params map[string]interface{}) bool

	Evaluate(expr string, params map[string]interface{}) interface{}
}

type SmartExpression struct {

}

func (s *SmartExpression) Eval(expr string, params map[string]interface{}) bool {
	if r := s.Evaluate(expr, params); r == nil {
		return false
	} else {
		return r.(bool)
	}
}

func (s *SmartExpression) Evaluate(expr string, params map[string]interface{}) interface{} {
	if expression, err := govaluate.NewEvaluableExpression(expr); err != nil {
		logs.Error("create expression failed. expr: %s, params: %v, err: %v", expr, params, err)
	} else if result, err := expression.Evaluate(params); err != nil {
		logs.Error("exec expression failed. expr: %s, params:%v, err: %v", expr, params, err)
	} else {
		return result
	}
	return nil
}


