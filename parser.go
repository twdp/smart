package smart

import (
	"errors"
	"fmt"
	"github.com/clbanning/mxj"
	"github.com/emirpasic/gods/lists/arraylist"
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
	Parse(element map[string]interface{}) (*NodeModel, error)
}

type ModelGen interface {
	newModel() *NodeModel
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
	container["end"] = &StartParserFactory{}

	return &DefaultSnakerParserContainer{
		container,
	}
}


type AbstractNodeParser struct {
	//model *model.NodeModel
	Parent ModelGen
}

func (a *AbstractNodeParser) Parse(element map[string]interface{}) (*NodeModel, error) {
	m := a.Parent.newModel()
	//a.model = model
	m.Name = element[AttrName].(string)
	m.DisplayName = element[AttrDisplayName].(string)
	// interceptor

	v := element[NodeTransition]
	tms := arraylist.New()

	if  v != nil {
		vv := reflect.ValueOf(v)

		switch vv.Kind() {
		case reflect.Map:
			tms.Add(v)
		case reflect.Slice:
			for _, k := range v.([]interface{}) {
				tms.Add(k)
			}
		}
	}


	for _,  te := range tms.Values() {
		tte := te.(map[string]interface{})
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
		m.Outputs.Add(transition)
	}

	a.parseNode(m, element)

	return m, nil
}

// 子类可覆盖此方法，完成特定的解析
func (a *AbstractNodeParser) parseNode(model *NodeModel, element map[string]interface{}) error {
	return nil
}

func (a *AbstractNodeParser) newModel() *NodeModel {
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

func (s *StartParser) newModel() *NodeModel {
	newNode := NewNodeModel("", "")
	startModel := &StartModel{ NodeModel: *newNode }
	newNode.Child = startModel


	return newNode
	//return (*NodeModel) (unsafe.Pointer(&StartModel{
	//	NodeModel: *newNode,
	//}))
}

///////////////////////////////////////////////////////////////////////////////////////////////////////

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
					process.Nodes.Add(m)
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
						process.Nodes.Add(m)
					}
				}
			}
		}

		for _, node := range process.Nodes.Values() {
			nodeModel := node.(*NodeModel)
			for _, t := range nodeModel.Outputs.Values() {
				transition := t.(*TransitionModel)
				to := transition.To
				for _, node2 := range process.Nodes.Values() {
					nodeModel2 := node2.(*NodeModel)
					if to == nodeModel2.Name {
						nodeModel2.Inputs.Add(transition)
						transition.Target = nodeModel2
					}
				}
			}
		}

		return process, nil
	}
	return nil, fmt.Errorf("解析xml失败")
}

// 对流程定义xml的节点，根据其节点对应的解析器解析节点内容
func (x *XmlParser) parseModel(node map[string]interface{}) (*NodeModel, error) {
	return x.ElementParserContainer.GetNodeParserFactory(node[ElementType].(string)).NewParse().Parse(node)
}
