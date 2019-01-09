package smart

import (
	"fmt"
	"github.com/clbanning/mxj"
	"github.com/emirpasic/gods/lists/arraylist"
	"reflect"
	"testing"
)

type AParserFactory struct {


}

func (a *AParserFactory) NewParse() NodeParser {
	return &AParser{}
}

type AParser struct {
	NodeParser
}



func (a *AParser) Parse(element map[string]interface{}) (*NodeModel, error) {
	return &NodeModel{
		BaseModel: BaseModel{
			Name: element[AttrName].(string),
			DisplayName: element[AttrDisplayName].(string),
		},
		Inputs:arraylist.New(),
		Outputs:arraylist.New(),
	}, nil
}

func TestXmlParser_ParseXml(t *testing.T) {
	x := `<process name="t" displayName="tt">
		<a name="test" displayName="测试一下">
		<b/>
	</a></process>`
	parse := &XmlParser{
		ElementParserContainer: NewDefaultSnakerParserContainer(),
	}
	parse.ElementParserContainer.AddParserFactory("a", &AParserFactory{})

	model, err := parse.ParseXml(x)
	fmt.Println(model, err)
	fmt.Println(model.Nodes.Size())

}


type MockNodeParserFactory struct {

}

type MockNodeParse struct {

}

func (m *MockNodeParse) Parse(element map[string]interface{}) (*NodeModel, error) {
	return &NodeModel{}, nil
}

func (m *MockNodeParserFactory) NewParse() NodeParser {
	return &MockNodeParse{}
}

func TestDefaultSnakerParserContainer_GetNodeParserFactory(t *testing.T) {
	c := NewDefaultSnakerParserContainer()
	m := &MockNodeParserFactory{}
	c.container["a"] = m

	defer func() {
		if err := recover(); err != nil {
			if err != "[b]没有对应的解析工厂类" {
				t.Error("SnakerParserContainer 测试失败")
			}
		}
	}()

	if c.GetNodeParserFactory("a") != m {
		t.Error("获取失败")
	}

	c.GetNodeParserFactory("b")
}

func TestParseXml(t *testing.T) {
	x := `<?xml version="1.0" encoding="UTF-8" standalone="no"?>

<process displayName="借款申请流程" instanceUrl="/snaker/flow/all" name="borrow">
    <start displayName="start1" layout="42,118,-1,-1" name="start1">
        <transition g="" name="transition1" offset="0,0" to="apply"/>
    </start>
    <end displayName="end1" layout="479,118,-1,-1" name="end1"/>
    <task assignee="apply.operator" autoExecute="Y" displayName="借款申请" form="/flow/borrow/apply" layout="126,116,-1,-1" name="apply" performType="ANY" taskType="Major">
        <transition g="" name="transition2" offset="0,0" to="approval"/>
    </task>
    <task assignee="approval.operator" autoExecute="Y" displayName="审批" form="/snaker/flow/approval" layout="252,116,-1,-1" name="approval" performType="ANY" taskType="Major">
     <transition g="" name="transition3" offset="0,0" to="decision1"/>
    </task>
    <decision displayName="decision1" expr="#result" layout="384,118,-1,-1" name="decision1">
        <transition displayName="同意" g="" name="agree" offset="0,0" to="end1"/>
        <transition displayName="不同意" g="408,68;172,68" name="disagree" offset="0,0" to="apply"/>
    </decision>

	<decision displayName="decision1" expr="#result" layout="384,118,-1,-1" name="decision1">
        <transition displayName="同意" g="" name="agree" offset="0,0" to="end1"/>
        <transition displayName="不同意" g="408,68;172,68" name="disagree" offset="0,0" to="apply"/>
    </decision>
</process>
`
	c, _ := mxj.NewMapXml([]byte(x))
	root := c.Old()
	vv := reflect.ValueOf(root["process"].(map[string]interface{})["decision"])
	switch vv.Kind() {
	case reflect.Map:
		fmt.Println("map")
	case reflect.Slice:
		fmt.Println("slice")
	}
	//fmt.Println(root["process"].(map[string]interface{})["-name"])

}

func TestStartParserFactory_NewParse(t *testing.T) {
	startParser := &StartParser{}
	startParser.AbstractNodeParser.Parent = startParser

	fmt.Println(startParser.Parse(map[string]interface{}{
		"-name": "dd",
		"-displayName": "fff",
	}))
}