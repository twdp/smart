package smart

import (
	"errors"
	"fmt"
	"github.com/clbanning/mxj"
	"reflect"
)

// xml 节点信息
const (
	RootElement = "process"
	// map中节点类型
	ElementType = "elementType"
	// 变迁节点名称
	NodeTransition = "transition"

	// 节点属性名称
	AttrName = "-name"
	AttrDisplayName = "-displayName"
	AttrInstanceUrl = "-instanceUrl"
	AttrInstanceNoClazz = "-instanceNoClass"
	AttrExpr = "-expr"
	AttrHandleClazz = "-handleClass"
	AttrForm = "-form"
	AttrField = "-field"
	AttrValue = "-value"
	AttrAttr = "-attr"
	attrType= "-type"
	AttrAssignee = "-assignee"
	AttrAssignmentHandler = "-assignmentHandler"
	AttrPerormType = "-performType"
	AttrTaskType = "-taskType"
	AttrTo = "-to"
	AttrProcessName = "-processName"
	AttrVersion = "-version"
	AttrExpireTime = "-expireTime"
	AttrAutoExecute = "-autoExecute"
	AttrCallback = "-callback"
	AttrReminderTime = "-reminderTime"
	AttrReminderRepeat = "-reminderRepeat"
	AttrClazz = "-clazz"
	AttrMethodName = "-methodName"
	AttrArgs = "-args"
	AttrVar = "-var"
	AttrLayout = "-layout"
	AttrG = "-g"
	AttrOffset = "-offset"
	AttrPreInterceptors = "-preInterceptors"
	AttrPostInterceptors = "-postInterceptors"
)

//
type SmartrParserContainer interface {

	// 添加解析
	AddParserFactory(elementName string, f NodeParserFactory)

	// 根据element名称获取对应的工厂
	GetNodeParserFactory(elementName string) NodeParserFactory
}

// engine上挂载factory
type NodeParserFactory interface {

	// 根据elementName查找使用哪个parser
	NewParse() NodeParser
}

// 节点解析接口
type NodeParser interface {

	// 节点dom元素解析方法，由实现类完成解析
	Parse(element map[string]interface{}) (INodeModel, error)
}

type ModelGen interface {
	newModel() INodeModel
}

type DefaultSnakerParserContainer struct {
	container map[string]NodeParserFactory
}

func (d *DefaultSnakerParserContainer) AddParserFactory(elementName string, f NodeParserFactory) {
	d.container[elementName] = f
}

func (d *DefaultSnakerParserContainer) GetNodeParserFactory(elementName string) NodeParserFactory {
	if f, ok := d.container[elementName]; ok {
		return f
	} else {
		panic(fmt.Sprintf("[%s]没有对应的解析工厂类", elementName))
	}
}

func NewDefaultSmartParserContainer() *DefaultSnakerParserContainer {
	container := make(map[string]NodeParserFactory)

	// 注册解析工厂
	container["start"] = &StartParserFactory{}
	container["end"] = &EndParserFactory{}
	container["task"] = &TaskParserFactory{}
	container["decision"] = &DecisionParserFactory{}
	container["custom"] = &CustomParserFactory{}


	return &DefaultSnakerParserContainer{
		container,
	}
}


type AbstractNodeParser struct {
	//model *model.NodeModel
	Parent ModelGen
}

func (a *AbstractNodeParser) Parse(element map[string]interface{}) (INodeModel, error) {
	m := a.Parent.newModel()

	m.SetName(element[AttrName].(string))
	m.SetDisplayName(element[AttrDisplayName].(string))

	// interceptor

	v := element[NodeTransition]
	var tms []map[string]interface{}
	if  v != nil {
		vv := reflect.ValueOf(v)

		switch vv.Kind() {
		case reflect.Map:
			tms = append(tms, v.(map[string]interface{}))
		case reflect.Slice:
			for _, k := range v.([]interface{}) {
				tms = append(tms, k.(map[string]interface{}))
			}
		}
	}


	for _,  tte := range tms {
		if _, ok := tte[AttrExpr]; !ok {
			tte[AttrExpr] = ""
		}
		if _, ok := tte[AttrDisplayName]; !ok {
			tte[AttrDisplayName] = ""
		}
		transition := &TransitionModel{
			BaseModel: BaseModel{
				Name: tte[AttrName].(string),
				DisplayName: tte[AttrDisplayName].(string),
			},
			To: tte[AttrTo].(string),
			Expr: tte[AttrExpr].(string),
			Source: m,
		}
		m.SetOutputs(append(m.GetOutputs(), transition))
		//m.Outputs = append(m.Outputs, transition)
	}

	//a.parseNode(m, element)

	if p, ok := a.Parent.(ParseNode); ok {
		if err :=  p.parseNode(m, element); err != nil {
			return nil, err
		}
	}
	return m, nil
}

type ParseNode interface {
	parseNode(model INodeModel, element map[string]interface{}) error
}

// 子类可覆盖此方法，完成特定的解析
//func (a *AbstractNodeParser) parseNode(model *NodeModel, element map[string]interface{}) error {
//	if p, ok := a.Parent.(ParseNode); ok {
//		return p.parseNode(model, element)
//	}
//	return nil
//}

func (a *AbstractNodeParser) newModel() interface{} {
	panic("未实现此方法")
}

///////////////////////////////////////////////////////////////////////////////////////////////////////

type StartParser struct {
	AbstractNodeParser
}

type StartParserFactory struct {

}

func (s *StartParserFactory) NewParse() NodeParser {
	ss := new(StartParser)
	ss.Parent = ss
	return ss
}

func (s *StartParser) newModel() INodeModel {
	newNode := NewNodeModel("", "")
	startModel := &StartModel{ NodeModel: *newNode }

	startModel.SetExec(startModel.exec)
	return startModel
}

///////////////////////////////////////////////////////////////////////////////////////////////////////

type EndParser struct {
	AbstractNodeParser
}

type EndParserFactory struct {

}

func (e *EndParserFactory) NewParse() NodeParser {
	end := new(EndParser)

	end.Parent = end
	return end
}

func (e *EndParser) newModel() INodeModel {
	newNode := NewNodeModel("", "")
	endModel := &EndModel{ NodeModel: *newNode }

	endModel.SetExec(endModel.exec)
	return endModel
}

//////////////////////////////////////////////////////////////////////////////////////////////////////

type CustomParser struct {
	AbstractNodeParser
}

type CustomParserFactory struct {

}

func (c *CustomParserFactory) NewParse() NodeParser {
	custom := new(CustomParser)
	custom.Parent = custom
	return custom
}

func (c *CustomParser) newModel() INodeModel {
	newNode := NewNodeModel("", "")
	workModel := WorkModel{ *newNode, "" }
	customModel := &CustomModel{ workModel, "" }

	customModel.SetExec(customModel.exec)
	return customModel
}

func (a *CustomParser) parseNode(model INodeModel, element map[string]interface{}) error {
	customModel := model.(*CustomModel)
	if element[AttrClazz] == nil {
		return errors.New("自定义模型需要指定clazz")
	}
	customModel.Clazz = element[AttrClazz].(string)

	return nil
}

//////////////////////////////////////////////////////////////////////////////////////////////////////

type DecisionParser struct {
	AbstractNodeParser
}

type DecisionParserFactory struct {

}

func (df *DecisionParserFactory) NewParse() NodeParser {
	d := new(DecisionParser)
	d.Parent = d
	return d
}

func (df *DecisionParser) newModel() INodeModel {
	newNode := NewNodeModel("", "")
	decisionModel := &DecisionModel{ NodeModel: *newNode }

	decisionModel.SetExec(decisionModel.exec)
	return decisionModel
}

//////////////////////////////////////////////////////////////////////////////////////////////////////

type TaskParser struct {
	AbstractNodeParser
}

type TaskParserFactory struct {

}

func (tf *TaskParserFactory) NewParse() NodeParser {
	t := new(TaskParser)
	t.Parent = t
	return t
}

func (tp *TaskParser) newModel() INodeModel {
	newNode := NewNodeModel("", "")
	workModel := WorkModel{ *newNode, "" }

	taskModel := &TaskModel{ WorkModel: workModel }

	taskModel.SetExec(taskModel.exec)
	return taskModel
}

func (tp *TaskParser) parseNode(model INodeModel, element map[string]interface{}) error {
	task := model.(*TaskModel)
	task.AssignTo = element[AttrAssignee].(string)

	return nil
}

//////////////////////////////////////////////////////////////////////////////////////////////////////

type Parser interface {
	ParseXml(content string) (*ProcessModel, error)
}


type XmlParser struct {
	// xml 元素解析容器
	ElementParserContainer SmartrParserContainer
}

// 解析流程定义文件，并将解析后的对象放入模型容器中
func (x *XmlParser) ParseXml(content string) (*ProcessModel, error) {
	if c, err := mxj.NewMapXml([]byte(content)); err != nil {
		return nil, errors.New(fmt.Sprintf("解析xml文件出错, content: %s", content))
	} else {
		// 根元素
		root := c.Old()[RootElement].(map[string]interface{})
		process := NewProcess(root[AttrName].(string), root[AttrDisplayName].(string))

		for k, v := range root {
			vv := reflect.ValueOf(v)
			switch vv.Kind() {
			case reflect.Map:
				vvv := v.(map[string]interface{})
				vvv[ElementType] = k
				if m, err := x.parseModel(vvv); err != nil {
					return nil, err
				} else {
					process.Nodes = append(process.Nodes, m)
				}
			case reflect.Slice:
				// 节点类型多个时
				// 是slice类型
				for _, kk := range v.([]interface{}) {
					vvv := kk.(map[string]interface{})
					vvv[ElementType] = k
					if m, err := x.parseModel(vvv); err != nil {
						return nil, err
					} else {
						process.Nodes = append(process.Nodes, m)
					}
				}
			}
		}

		for _, node := range process.Nodes {
			for _, transition := range node.GetOutputs() {
				to := transition.To
				for _, node2 := range process.Nodes {
					if to == node2.GetName() {
						transition.Target = node2
						node2.SetInputs(append(node2.GetInputs(), transition))
					}
				}
				if transition.Target == nil {
					panic("transition target is nil. ---> " + transition.Name + " ---> targetName: " + transition.To)
				}
			}
		}

		return process, nil
	}
	return nil, fmt.Errorf("解析xml失败")
}

// 对流程定义xml的节点，根据其节点对应的解析器解析节点内容
func (x *XmlParser) parseModel(node map[string]interface{}) (INodeModel, error) {
	return x.ElementParserContainer.GetNodeParserFactory(node[ElementType].(string)).NewParse().Parse(node)
}
