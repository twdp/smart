package smart

import (
	"errors"
	"fmt"
	"sync"
)

type ProcessService interface {

	// 检查流程定义对象
	Check(process *Process, idOrName string) error

	// 根据主键ID获取流程定义对象
	GetProcessById(id int64) *Process

	ParseProcess(process *Process) (*ProcessModel, error)
}

type ProcessAccess interface {

}

type SmartProcessService struct {
	sync.RWMutex
	engine Engine

	Child ProcessAccess
}

func NewSmartProcessService(engine Engine) ProcessService {
	return &SmartProcessService{
		engine: engine,
	}
}

func (s *SmartProcessService) Check(process *Process, idOrName string) error {
	if nil == process {
		return errors.New(fmt.Sprintf("指定的流程定义[id/name=%s]不存在", idOrName))
	} else if process.Status == ProcessInit {
		return errors.New(fmt.Sprintf("指定的流程定义[id/name=%s,version=%d]为非活动状态", idOrName, process.Version))
	}
	return nil
}

func (s *SmartProcessService) ParseProcess(process *Process) (*ProcessModel, error) {
	return s.engine.Parser().ParseXml(process.Content)
}


func (s *SmartProcessService) GetProcessById(id int64) *Process {
	return &Process{
		Id: 1,
		Content: `<?xml version="1.0" encoding="UTF-8" standalone="no"?>

<process displayName="借款申请流程" instanceUrl="/snaker/flow/all" name="borrow">
    <start displayName="start1" layout="42,118,-1,-1" name="start1">
    </start>
    <end displayName="end1" layout="479,118,-1,-1" name="end1"/>
    
    
</process>`,
	}
}