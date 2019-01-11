package smart

import (
	"fmt"
	"github.com/astaxie/beego/logs"
	"github.com/pkg/errors"
	"strconv"
	"tianwei.pro/kit/di"
)

var Di = di.New()

type Engine interface {

	Parser() Parser

	Cache() CacheManager

	Instance() InstanceService

	// 获取流程处理service
	Process() ProcessService

	// 获取表达式引擎
	Expression() Expression

	// 根据流程定义ID启动流程实例
	StartInstanceById(id int64) (*Instance, error)

	// 根据流程定义id和操作人|flag启动流程实例
	StartInstanceByIdAndOperator(id int64, operator string) (*Instance, error)

	// 根据流程定义id和操作人|flag和参数启动流程实例
	StartInstanceByIdAndOperatorAndArgs(id int64, operator string, args map[string]interface{}) (*Instance, error)

}

// smart engine
type SmartEngine struct {

	instance InstanceService

	// 流程处理service
	process ProcessService

	// 表达式引擎
	expression Expression

	// 缓存控制器
	cache CacheManager

	parser Parser
}

func NewSmartEngine() Engine {
	engine := &SmartEngine{}

	i := NewSmartInstanceService(engine)
	p := NewSmartProcessService(engine)
	e := NewSmartExpression()
	c := NewSmartCacheManager()
	x := &XmlParser{
		NewDefaultSmartParserContainer(),
	}

	engine.instance = i
	engine.parser = x
	engine.expression = e
	engine.process = p
	engine.cache = c

	return engine
}

func (s *SmartEngine) Parser() Parser {
	if s.parser == nil {
		panic("流程解析实例未设置")
	}
	return s.parser
}

func (s *SmartEngine) Cache() CacheManager {
	if s.cache == nil {
		panic("缓存管理器未设置")
	}
	return s.cache
}

func (s *SmartEngine) Instance() InstanceService {
	if s.instance == nil {
		panic("未设置流程实例service")
	}
	return s.instance
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

func (s *SmartEngine) StartInstanceById(id int64) (*Instance, error) {
	return s.StartInstanceByIdAndOperatorAndArgs(id, "", nil)
}

func (s *SmartEngine) StartInstanceByIdAndOperator(id int64, operator string) (*Instance, error) {
	return s.StartInstanceByIdAndOperatorAndArgs(id, operator, nil)
}

func (s *SmartEngine) StartInstanceByIdAndOperatorAndArgs(id int64, operator string, args map[string]interface{}) (*Instance, error) {
	if nil == args {
		args = make(map[string]interface{})
	}
	process := s.Process().GetProcessById(id)
	if err := s.Process().Check(process, strconv.FormatInt(id, 10)); err != nil {
		return nil, err
	}
	if process.Status != ProcessRunning {
		logs.Error("process not running. id: %d", id)
		return nil, errors.New("流程未激活")
	}
	return s.startProcess(process, operator, args)
}

func (s *SmartEngine) startProcess(process *Process, operator string, args map[string]interface{}) (*Instance, error) {
	if context, err := s.execute(process, operator, args, 0, ""); nil != err {
		return nil, err
	} else if pm, err := s.Process().ParseProcess(process); err != nil {
		return nil, err
	} else {
		context.ProcessModel = pm
		if pm != nil {
			if start, err := pm.GetStart(); err != nil {
				return nil, err
			} else if err = start.Execute(context); err != nil {
				return nil, err
			}
		}

		return context.Instance, nil
	}
}

func (s *SmartEngine) execute(process *Process, operator string, args map[string]interface{}, parentId int64, parentNodeName string) (*Context, error) {
	if instance, err := s.Instance().CreateInstanceUseParentInfo(process, operator, args, parentId, parentNodeName); err != nil {
		return nil, err
	} else {
		fmt.Println(instance)
		current := &Context{
			Engine: s,
			Instance: instance,
			Args: args,
			Operator: operator,
		}

		return current, nil
	}

}