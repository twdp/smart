package smart

import (
	"fmt"
	"testing"
)

func TestSmartExpression_Eval(t *testing.T) {
	s := SmartExpression{}


	parameters := make(map[string]interface{}, 8)
	parameters["foo"] = -1

	if s.Eval("foo > 0", parameters) {
		t.Errorf("执行失败啦")
	}

	parameters["foo"] = 1
	fmt.Println(s.Eval("foo > 0", parameters))

}