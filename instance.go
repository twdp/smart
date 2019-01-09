package smart

type InstanceService interface {

	CreateInstanceUseParentInfo(process *Process, operator string, args map[string]interface{}, parentId int64, parentNodeName string) (*Instance, error)
}

type SmartInstanceService struct {
	Engine Engine
}

func NewSmartInstanceService(engine Engine) InstanceService {
	return &SmartInstanceService{
		Engine: engine,
	}
}

func (s *SmartInstanceService) CreateInstanceUseParentInfo(process *Process, operator string, args map[string]interface{}, parentId int64, parentNodeName string) (*Instance, error) {
	return &Instance{}, nil
}