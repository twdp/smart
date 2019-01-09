package smart

import "tianwei.pro/kit/di"

var Di = di.New()

type Engine interface {

	// 获取流程处理service
	Process() ProcessService

	// 获取表达式引擎
	Expression() Expression

	// 部署流程
	Deploy(process *Process)
}

// smart engine
type SmartEngine struct {

	// 流程处理service
	process ProcessService

	// 表达式引擎
	expression Expression

	// 缓存控制器
	cache CacheManager
}

func NewSmartEngine() Engine {
	engine := &SmartEngine{}

	p := NewSmartProcessService(engine)
	e := NewSmartExpression()
	c := NewSmartCacheManager()

	engine.expression = e
	engine.process = p
	engine.cache = c

	return engine
}

// 获取流程处理service
func (s *SmartEngine) Process() ProcessService {
	if s.process == nil {
		panic("未设置流程处理service")
	}
	return s.process
}

// 获取表达式引擎
func (s *SmartEngine) Expression() Expression {
	if s.expression == nil {
		panic("未设置解析引擎")
	}
	return s.expression
}

// 部署流程
func (s *SmartEngine) Deploy(process *Process) {

	panic("implement me")
}
