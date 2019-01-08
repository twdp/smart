package smart

import "tianwei.pro/kit/di"

var Di = di.New()

type Engine interface {

	// 获取表达式引擎
	Expression() Expression
}