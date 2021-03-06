package smart

import (
	"github.com/astaxie/beego/logs"
	"github.com/astaxie/beego/orm"
	"github.com/pkg/errors"
)

type InstanceService interface {

	// 根据流程、操作人员、父流程实例ID创建流程实例
	CreateInstanceUseParentInfo(process *Process, operator string, args map[string]interface{}, parentId int64, parentNodeName string) (*Instance, error)

	// 更新流程实例
	UpdateInstance(instance *Instance) error
}

type SmartInstanceService struct {
	Engine Engine
}

func NewSmartInstanceService(engine Engine) InstanceService {
	return &SmartInstanceService{
		Engine: engine,
	}
}


 // 根据流程、操作人员、父流程实例ID创建流程实例
 // @param process 流程定义对象
 // @param operator 操作人员ID
 // @param args 参数列表
 // @param parentId 父流程实例ID
 // @param parentNodeName 父流程节点模型
 // @return 活动流程实例对象
func (s *SmartInstanceService) CreateInstanceUseParentInfo(process *Process, operator string, args map[string]interface{}, parentId int64, parentNodeName string) (*Instance, error) {
	instance := &Instance{
		Name: process.Name,
		DisplayName: process.DisplayName,
		ProcessId: process.Id,
		Content: process.Content,
		Deployer: operator,
		ParentId: parentId,
		ParentNodeName: parentNodeName,
	}
	instance.SetVariable(args)
	if _, err := orm.NewOrm().Insert(instance); err != nil {
		logs.Error("create instance failed. process: %v, err: %v", *process, err)
		return nil, errors.New("创建流程实例失败")
	}

	return instance, nil
}

func (s *SmartInstanceService) UpdateInstance(instance *Instance) error {
	_, err := orm.NewOrm().Update(instance)
	if nil != err {
		logs.Error("update instance failed. instance: %v, err: %v", instance, err)
		return errors.New("更新实例失败")
	}
	return nil
}