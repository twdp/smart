package smart

import (
	"errors"
	"fmt"
	"github.com/astaxie/beego/logs"
	"github.com/astaxie/beego/orm"
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
	p := &Process{ Id: id }
	err := orm.NewOrm().Read(p)
	if err != nil {
		logs.Error("read process failed. id: %d, err: %v", id, err)
		return nil
	}
	return p
//	return &Process{
//		Id: id,
//		Content: `<?xml version="1.0" encoding="UTF-8" standalone="no"?>
//
//<process displayName="借款申请流程" instanceUrl="/snaker/flow/all" name="borrow">
//    <start displayName="start1" layout="42,118,-1,-1" name="start1">
//        <transition g="" displayName="xx" name="transition1" offset="0,0" to="apply"/>
//    </start>
//    <end displayName="end1" layout="479,118,-1,-1" name="end1"/>
//
//    <task assignee="apply.operator" autoExecute="Y" displayName="借款申请" form="/flow/borrow/apply" layout="126,116,-1,-1" name="apply" performType="ANY" taskType="Major">
//        <transition g="" displayName="xx" name="transition2" offset="0,0" to="approval"/>
//    </task>
//    <task assignee="approval.operator" autoExecute="Y" displayName="审批" form="/snaker/flow/approval" layout="252,116,-1,-1" name="approval" performType="ANY" taskType="Major">
//     <transition g="" displayName="xx" name="transition3" offset="0,0" to="decision1"/>
//    </task>
//
//</process>`,
//	}
}