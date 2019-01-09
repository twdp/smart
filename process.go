package smart

import (
	"sync"
)

type ProcessService interface {

	// 部署流程实例
	Deploy(process *Process) error

	// 将制定id的流程状态设置为可用
	// 其他同名称的流程均置为停止
	Start(id int64) error
}

type ProcessAccess interface {

}

type SmartProcessService struct {
	sync.RWMutex
	engine Engine

	Child ProcessAccess
}

func (s *SmartProcessService) Start(id int64) error {
	return nil

}

func NewSmartProcessService(engine Engine) ProcessService {
	return &SmartProcessService{
		engine: engine,
	}
}

func (s *SmartProcessService) Deploy(process *Process) error {

	return nil

}
