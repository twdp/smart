package smart

import (
	"container/list"
	"errors"
	"fmt"
	"github.com/astaxie/beego/logs"
	"github.com/emirpasic/gods/lists"
	"github.com/emirpasic/gods/lists/arraylist"
	"reflect"
)

type BaseModel struct {

	// 元素名称
	Name string

	// 显示名称
	DisplayName string
}

// 将执行对象execution交给具体的处理器处理
func (b *BaseModel) fire(handler Handler, ctx *Context) error {
	return handler.Handle(ctx)
}

//////////////////////////////////////////////////////////////////////////////

type NodeModel struct {
	BaseModel

	Inputs lists.List

	Outputs lists.List

	// 前置局部拦截器实例集合
	PreInterceptors lists.List

	// 后置局部拦截器实例集合
	PostInterceptors lists.List

	Child NodeModelChild
}

//
type NodeModelChild interface {
	exec(context *Context) error
}

func NewNodeModel(name, displayName string) *NodeModel {
	return &NodeModel{
		BaseModel: BaseModel{
			Name: name,
			DisplayName: displayName,
		},
		Inputs: arraylist.New(),
		Outputs: arraylist.New(),
		PreInterceptors: arraylist.New(),
		PostInterceptors: arraylist.New(),
	}
}

//  对执行逻辑增加前置、后置拦截处理
func (n *NodeModel) Execute(context *Context) error {
	if err := n.intercept(n.PreInterceptors, context); err != nil {
		return err
	} else if err = n.exec(context); err != nil {
		return err
	} else if err = n.intercept(n.PostInterceptors, context); err != nil {
		return err
	}
	return nil
}

// 具体节点模型需要完成的执行逻辑
func (n *NodeModel) exec(context *Context) error {
	if n.Child == nil {
		panic("初始化NodeModel时，请设置Child")
	}
	return n.Child.exec(context)
}

// 拦截方法
func (n *NodeModel) intercept(interceptors lists.List, context *Context) error {
	for _, v := range interceptors.Values() {
		interceptor := v.(Interceptor)
		if err := interceptor.Intercept(context); err != nil {
			return err
		}
	}
	return nil
}

// 运行变迁继续执行
func (n *NodeModel) runOutTransition(context *Context) error {
	for _, v := range n.Outputs.Values() {
		tm := v.(*TransitionModel)
		tm.Enable = true
		if err := tm.Execute(context); err != nil {
			return err
		}
	}
	return nil
}


/**
 * 根据父节点模型、当前节点模型判断是否可退回。可退回条件：
 * 1、满足中间无fork、join、subprocess模型
 * 2、满足父节点模型如果为任务模型时，参与类型为any
 */
func (n *NodeModel) CanRejected(current *NodeModel, parent *NodeModel) bool {
	switch t := (interface{})(parent).(type) {
	case *TaskModel:
		return t.PerformType == PerformtypeAll
	}
	result := false
	for _, e := range n.Outputs.Values() {
		tm := e.(*TransitionModel)
		source := tm.Source
		if source == parent {
			return true
		}
		switch s := (interface{})(source).(type) {
		case *ForkModel:
			logs.Debug("can rejected source. %v", s)
			continue
		case *JoinModel:
			logs.Debug("can rejected source. %v", s)
			continue
		case *SubProcessModel:
			logs.Debug("can rejected source. %v", s)
			continue
		case *StartModel:
			logs.Debug("can rejected source. %v", s)
			continue
		}
		result = result || n.CanRejected(source, parent)
	}
	return result
}

func (n *NodeModel) getNextModels(clazz interface{}) lists.List {
	r := arraylist.New()
	c := reflect.TypeOf(clazz)
	for _, o := range n.Outputs.Values() {
		n.AddNextModels(r, o.(*TransitionModel), c)
	}
	return r
}

func (n *NodeModel) AddNextModels(r lists.List, tm *TransitionModel, t reflect.Type) {
	target := reflect.TypeOf(tm.Target)
	if t.AssignableTo(target) {
		r.Add(tm.Target)
	} else {
		for _, o := range tm.Target.Outputs.Values() {
			n.AddNextModels(r, o.(*TransitionModel), t)
		}
	}
}


//////////////////////////////////////////////////////////////////////////////////////

type TransitionModel struct {
	BaseModel

	// 当前转移路径是否可用
	Enable bool

	// 变迁的目标节点应用
	Target *NodeModel

	// 变迁的源节点引用
	Source *NodeModel

	// 变迁的目标节点name名称
	To string

	//  变迁的条件表达式，用于decision
	Expr string

	// 转折点图形数据
	// G string
}

func (t *TransitionModel) Execute(context *Context) error {
	if !t.Enable {
		return nil
	}
	//如果目标节点模型为TaskModel，则创建task
	if isTask, ok := t.Target.Child.(*TaskModel); ok {
		return t.fire(&CreateTaskHandler{
			TaskModel: isTask,
		}, context)
	} else if isSubProcess, ok := t.Target.Child.(*SubProcessModel); ok {
		//如果目标节点模型为SubProcessModel，则启动子流程

		return t.fire(&StartSubProcessHandler{
			SubProcessModel: isSubProcess,
		}, context)
	} else {
		//如果目标节点模型为其它控制类型，则继续由目标节点执行
		return t.Target.Execute(context)
	}
	return nil
}

// 开始节点定义start元素
type StartModel struct {
	NodeModel
}

func (s *StartModel) exec(context *Context) error {
	return s.runOutTransition(context)
}

// 结束节点end元素
type EndModel struct {

	NodeModel
}

func (e *EndModel) exec(context *Context) error {
	return e.fire(&EndProcessHandler{}, context)
}


// 工作元素
type WorkModel struct {
	NodeModel

	Form string
}

// 用户自定义处理
// snaker对外提供一个di容器
// 实现接口并注入到容器中
// snaker 处理时调用
type CustomModel struct {
	WorkModel

	// 实例名称
	Clazz string

	// 传入参数
	//Args string

}


// 从di容器中查找指定的实例
func (c *CustomModel) exec(context *Context) error {
	return Di.GetByName(c.Clazz).(Delegation).Execute(context)
}

// 决策定义decision元素
type DecisionModel struct {
	NodeModel

	// 决策选择表达式串（需要表达式引擎解析）
	Expr string
}


func (d *DecisionModel) exec(context *Context) error {
	logs.Info("%d->decision execution.getArgs():%v", context.Instance.Id, context.Args)

	isFound := false
	for _, e := range d.Outputs.Values() {
		tm := e.(*TransitionModel)

		if "" != tm.Expr && context.Engine.Expression().Eval(tm.Expr, context.Args) {
			tm.Enable = true
			tm.Execute(context)
			isFound = true
		}
	}

	if !isFound {
		return fmt.Errorf("%d->decision节点无法确定下一步执行路线", context.Instance.Id)
	}
	return nil
}

type ForkModel struct {
	NodeModel

}

func (f *ForkModel) exec(context *Context)error {
	return f.runOutTransition(context)
}

// 合并定义join元素
type JoinModel struct {
	NodeModel

}

func (j *JoinModel) exec(context *Context) error {
	if err := j.fire(&MergeBranchHandler{ JoinModel: j }, context); err != nil {
		return err
	} else if context.IsMerged {
		return j.runOutTransition(context)
	}
	return nil
}
// todo::
type ProcessModel struct {

	BaseModel

	// 节点元素集合
	Nodes lists.List

	TaskModels lists.List

	Process *Process
}

func NewProcess(name, displayName string) *ProcessModel {
	return &ProcessModel {
		BaseModel: BaseModel {
			Name: name,
			DisplayName: displayName,
		},
		Nodes: arraylist.New(),
		TaskModels: arraylist.New(),
	}
}


func (p *ProcessModel) GetWorkModels() list.List {
	r := list.New()
	for _, e := range p.Nodes.Values() {
		if _, ok := (e.(*NodeModel).Child).(*WorkModel); ok {
			r.PushBack(e)
		}
	}
	return *r
}

func (p *ProcessModel) GetStart() (*NodeModel, error) {
	for _, e := range p.Nodes.Values() {
		if _, ok := (e.(*NodeModel).Child).(*StartModel); ok {
			return e.(*NodeModel), nil
		}
	}
	return nil, errors.New("没有start节点")
}

func (p *ProcessModel) GetNode(nodeName string) (*NodeModel, error) {
	for _, e := range p.Nodes.Values() {
		if v, ok := e.(*NodeModel); ok {
			if v.Name == nodeName {
				return v, nil
			}
		}
	}

	return nil, errors.New(fmt.Sprintf("没有[%s]节点", nodeName))
}


type SubProcessModel struct {
	WorkModel

	ProcessName string

	Version int

	SubProcess *ProcessModel
}

func (s *SubProcessModel) exec(context *Context) error {
	return s.runOutTransition(context)
}


type TaskModel struct {
	WorkModel

	PerformType int8

	TaskType int8

	// 期望用时
	ExpectTime string

	// 提醒时间
	RemindTime string

	// 提醒间隔(分钟)
	RemindRepeat string

	// 是否自动执行
	AutoExecute bool
}

func (t *TaskModel) Exec(context *Context) error {
	//  any方式，直接执行输出变迁
	// all方式，需要判断是否已全部合并
	// 由于all方式分配任务，是每个执行体一个任务
	// 那么此时需要判断之前分配的所有任务都执行完成后，才可执行下一步，否则不处理
	if t.PerformType == PerformtypeAny {
		return t.runOutTransition(context)
	} else if err := t.fire(&MergeActorHandler{ TaskName: t.Name }, context); err != nil {
		return err
	} else if context.IsMerged {
		return t.runOutTransition(context)
	}
	return nil
}
