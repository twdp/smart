package smart


const (
	ProcessInit = iota
	ProcessRunning
	ProcessStop


	PerformtypeAll = iota

)

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
	Content string
}

/////////////////////////////////////////////////////////////////////////////////////////////////////

// 流程工作单实体类（一般称为流程实例）
type Instance struct {

	Id int64

	Name string

	DisplayName string

	ProcessId int64

	// 流程实例内容
	Content string

	// 发布者
	Deployer string


	// 流程实例附属变量
	Variable map[string]interface{}
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
	AssignTo string

	Variable map[string]interface{}
}

//////////////////////////////////////////////////////////////////////////////////////////////////////

