package smart

import (
	"errors"
	"fmt"
	"github.com/astaxie/beego/logs"
	"github.com/astaxie/beego/orm"
	"sync"
	"time"
)

const (
	processCache = "_processCache_"
	processCacheUseName = "_name_"
	processCacheUseId = "_id_"
)
type ProcessService interface {

	// 检查流程定义对象
	Check(process *Process, idOrName string) error

	// 根据主键ID获取流程定义对象
	GetProcessById(id int64) *Process

	ParseProcess(content string) (*ProcessModel, error)

	SaveProcess(process *Process) error

	ActiveProcess(id int64) error
}

type SmartProcessService struct {
	sync.RWMutex
	engine Engine
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

func (s *SmartProcessService) ParseProcess(content string) (*ProcessModel, error) {
	return s.engine.Parser().ParseXml(content)
}


func (s *SmartProcessService) GetProcessById(id int64) *Process {
	c := s.engine.Cache().Get(processCache)
	if p := c.Get(fmt.Sprintf("%s%d", processCacheUseId, id)); p != nil {
		return p.(*Process)
	}
	p := &Process{ Id: id }
	err := orm.NewOrm().Read(p)
	if err != nil {
		logs.Error("read process failed. id: %d, err: %v", id, err)
		return nil
	}
	c.Put(fmt.Sprintf("%s%d", processCacheUseId, id), p, 1 * time.Hour)
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

func (s *SmartProcessService) SaveProcess(process *Process) error {
	if _, err := s.ParseProcess(process.Content); err != nil {
		return err
	} else {
		s.Lock()
		defer s.Unlock()
		o := orm.NewOrm()
		p := &Process{
			Name: process.Name,
		}
		var olds []*Process
		_, err := o.QueryTable(p).Filter("name", process.Name).OrderBy("-version").All(&olds)
		if err != nil {
			logs.Error("query process by name failed. p: %v, err: %v", process, err)
			return errors.New("查询流程失败")
		}
		if len(olds) > 0 {
			process.Version = olds[0].Version + 1
		}
		if _, err = o.Insert(process); err != nil {
			logs.Error("insert process failed. process: %v, err: %v", process, err)
			return errors.New("创建流程失败")
		}
		return nil
	}
}

func (s *SmartProcessService) ActiveProcess(id int64) error {
	o := orm.NewOrm()
	p := &Process{ Id: id}
	if err := o.Read(p); err != nil {
		logs.Error("read process failed. id: %d, err: %v", id, err)
		return errors.New("查询流程失败")
	}
	o.Begin()
	_, err := o.QueryTable(p).Filter("name", p.Name).Filter("status", ProcessRunning).Update(orm.Params{
			"Status": ProcessStop,
			"UpdatedAt": time.Now(),
		})
	if err != nil {
		logs.Error("update process status failed. name: %s, err: %v", p.Name, err)
		return errors.New("更新流程失败")
	}
	p.Status = ProcessRunning
	if _, err = o.Update(&Process{ Id: id, Status: ProcessRunning}, "Status", "UpdatedAt"); err != nil {
		logs.Error("active process failed. id: %d, err: %v", id, err)
		o.Rollback()
		return errors.New("激活流程失败")
	} else {
		o.Commit()
	}
	s.engine.Cache().Get(processCache).Put(fmt.Sprintf("%s%s", processCacheUseName, p.Name), p, 1 * time.Hour)
	return nil
}