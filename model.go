package smart

import (
	"container/list"
	"errors"
	"fmt"
	"github.com/astaxie/beego/logs"
	"github.com/emirpasic/gods/lists"
	"github.com/emirpasic/gods/lists/arraylist"
)

type BaseModel struct {

	// 元素名称
	Name string

	// 显示名称
	DisplayName string
}

// todo: fire


//////////////////////////////////////////////////////////////////////////////

type NodeModel struct {
	BaseModel

	Inputs lists.List

	Outputs lists.List

	// 前置局部拦截器实例集合
	PreInterceptors lists.List

	// 后置局部拦截器实例集合
	PostInterceptors lists.List

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
	panic("子模型需要实现exec方法")
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
	return nil
}


type StartModel struct {
	NodeModel
}

func (s *StartModel) Execute(context *Context) error {
	return s.runOutTransition(context)
}


type EndModel struct {

	NodeModel
}

// todo::
func (e *EndModel) exec(context *Context) error {
	//e.Fire(, execution)
	return nil
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
	logs.Info("%d->decision execution.getArgs():%v", 11, context.Args)

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
		return errors.New(fmt.Sprintf("%d->decision节点无法确定下一步执行路线", 11))
	}
	return nil
}

type ForkModel struct {
	NodeModel

}

func (f *ForkModel) exec(context *Context)error {
	return f.runOutTransition(context)
}

type JoinModel struct {
	NodeModel

}

// todo::
//func (j *JoinModel) exec(execution *Execution) error {
//
//}


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
		if v, ok := e.(*WorkModel); ok {
			r.PushBack(v)
		}
	}
	return *r
}

func (p *ProcessModel) GetStart() (*StartModel, error) {
	for _, e := range p.Nodes.Values() {
		if v, ok := e.(*StartModel); ok {
			return v, nil
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

	PerformType int8

}

