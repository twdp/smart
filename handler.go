package smart

import "fmt"

// 流程各模型操作处理接口
type Handler interface {

	// 子类需要实现的方法，来处理具体的操作
	Handle(ctx *Context) error
}

// 结束流程实例的处理器
type EndProcessHandler struct {

}

// 结束当前流程实例，如果存在父流程，则触发父流程继续执行
func (e *EndProcessHandler) Handle(ctx *Context) error {
	// 查询实例的所有活跃的任务
	// 如果是主办任务，则跑异常
	// 否则将任务自动完成
	// 结束当前流程实例

	// 如果存在父流程，则重新构造Context执行对象，交给父流程的SubProcessModel模型execute
	fmt.Println("流程结束")
	return nil
}

type AbstractMergeHandler struct {
	Child MergeHandler
}

type MergeHandler interface {
	findActiveNodes() []string
}

func (a *AbstractMergeHandler) Handle(ctx *Context) error {
	panic("implement me")
}

type MergeActorHandler struct {
	AbstractMergeHandler

	TaskName string
}

type MergeBranchHandler struct {
	AbstractMergeHandler

	JoinModel *JoinModel
}

type CreateTaskHandler struct {
	TaskModel *TaskModel
}

func (c *CreateTaskHandler) Handle(ctx *Context) error {
	//panic("implement me")
	fmt.Println("处理任务啦....")
	return nil
}

type StartSubProcessHandler struct {
	SubProcessModel *SubProcessModel
}

func (s *StartSubProcessHandler) Handle(ctx *Context) error {
	panic("implement me")
}
