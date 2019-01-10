package smart

import (
	"encoding/json"
	"github.com/astaxie/beego/logs"
	"github.com/astaxie/beego/orm"
	"time"
)

const (
	ProcessInit = iota
	ProcessRunning
	ProcessStop


	PerformtypeAll = iota // 参与者fork任务,所有人均需要处理
	PerformtypeAny   // 普通任务

)

type Base struct {

	CreatedAt time.Time `orm:"auto_now_add"`

	UpdatedAt time.Time `orm:"auto_now"`
}

func init() {
	orm.RegisterModelWithPrefix("smart_", &Process{}, &Instance{}, &Task{})
}
//////////////////////////////////////////////////////////////////////////////////////////////////////

// 流程定义实体类
type Process struct {

	// 主键
	Id int64

	// 版本
	Version int

	// 流程定义名称，根据此字段启动流程
	Name string

	// 页面上展示的名称
	DisplayName string

	// 当前状态
	Status int8

	// 流程定义内容
	Content string `orm:"type(text)"`

	// todo:: 现在只针对TaskModel
	// 比如: 申请权限时，提交之后就可以看到整个流程
	// 是否预生成任务
	PreGeneratedTask bool
}

/////////////////////////////////////////////////////////////////////////////////////////////////////

// 流程工作单实体类（一般称为流程实例）
type Instance struct {

	Id int64

	Name string

	DisplayName string

	ProcessId int64

	// 流程实例内容
	Content string `orm:"type(text)"`

	// 发布者
	Deployer string `orm:"size(64);index"`// 谁创建的

	// 流程实例附属变量
	variable map[string]interface{} `orm:"-"`

	VariableJson string `orm:"type(text);default()"`

	// 创建子流程时，在哪个流程和节点上创建出来的
	ParentId int64

	ParentNodeName string

	Base

}

func (i *Instance) SetVariable(m map[string]interface{}) {
	mm, err := json.Marshal(m)
	if err != nil {
		logs.Error("marshal variable failed. m: %v, err: %v", m, err)
	}
	i.VariableJson = string(mm)
}

func (i *Instance) GetVariable() map[string]interface{} {
	if i.variable != nil {
		return i.variable
	}
	if i.variable == nil && i.VariableJson != ""{
		i.variable = make(map[string]interface{})
		err := json.Unmarshal([]byte(i.VariableJson), i.variable)
		if err != nil {
			logs.Error("unmarshal instance variable failed. variable: %v, err: %v", i.VariableJson, err)
		}
		return i.variable
	} else if i.variable == nil {
		i.variable = make(map[string]interface{})
		return i.variable
	} else {
		return i.variable
	}
}

func (i *Instance) AddVariable(n string, p interface{}) {
	i.GetVariable()[n] = p
	b, err := json.Marshal(i.variable)
	if err != nil {
		logs.Error("marshal variable failed. m: %v, err: %v", i.variable, err)
	}
	i.VariableJson = string(b)
}

//////////////////////////////////////////////////////////////////////////////////////////////////////

// task
// 提醒时间
// 多长时间自动过期
// 延迟执行时间
type Task struct {
	Id int64

	InstanceId int64

	TaskName string

	TaskDisplayName string

	// 任务分配给谁
	AssignTo string `orm:"size(64);index"`

	variable map[string]interface{} `orm:"-"`

	VariableJson string `orm:"type(text)"`

	Base

}


func (i *Task) SetVariable(m map[string]interface{}) {
	mm, err := json.Marshal(m)
	if err != nil {
		logs.Error("marshal variable failed. m: %v, err: %v", m, err)
	}
	i.VariableJson = string(mm)
}

func (i *Task) GetVariable() map[string]interface{} {
	if i.variable == nil && i.VariableJson != ""{
		i.variable = make(map[string]interface{})
		err := json.Unmarshal([]byte(i.VariableJson), i.variable)
		if err != nil {
			logs.Error("unmarshal instance variable failed. variable: %v, err: %v", i.VariableJson, err)
		}
		return i.variable
	} else if i.variable == nil {
		i.variable = make(map[string]interface{})
		return i.variable
	} else {
		return i.variable
	}
}

func (i *Task) AddVariable(n string, p interface{}) {
	i.GetVariable()[n] = p
}


//////////////////////////////////////////////////////////////////////////////////////////////////////

