package main

import (
	"fmt"
	"github.com/astaxie/beego/orm"
	"tianwei.pro/smart"
)
import _ "tianwei.pro/smart"
import _ "github.com/go-sql-driver/mysql"


func init() {
	//orm.DefaultTimeLoc = time.Local

	orm.RegisterDataBase("default", "mysql", "root:anywhere@tcp(127.0.0.1:3306)/smart?charset=utf8&loc=Asia%2FShanghai", 30)
	orm.RunSyncdb("default", false, true)

	smart.Di.Provide("load", &AA{})
}

type AA struct {

}

func (a *AA) Execute(context *smart.Context) error {
	context.Args["result"] = true
	fmt.Println("custom ....")
	return nil
}

func main() {

	orm.Debug = true

	e := smart.NewSmartEngine()

//	s := &smart.Process{
//		Name: "xx",
//		DisplayName: "ccc",
//		Content: `<?xml version="1.0" encoding="UTF-8" standalone="no"?>
//
//<process displayName="借款申请流程" instanceUrl="/snaker/flow/all" name="borrow">
//   <start displayName="start1" layout="42,118,-1,-1" name="start1">
//       <transition g="" name="transition1" offset="0,0" to="apply"/>
//   </start>
//   <end displayName="end1" layout="479,118,-1,-1" name="end1"/>
//
//   <task assignee="apply.operator" autoExecute="Y" displayName="借款申请" form="/flow/borrow/apply" layout="126,116,-1,-1" name="apply" performType="ANY" taskType="Major">
//       <transition g="" name="transition2" offset="0,0" to="approval"/>
//   </task>
//
//
//   <task assignee="approval.operator" autoExecute="Y" displayName="审批" form="/snaker/flow/approval" layout="252,116,-1,-1" name="approval" performType="ANY" taskType="Major">
//    <transition g="" name="transition3" offset="0,0" to="load"/>
//   </task>
//	<custom name="load" clazz="load" displayName="获取用户信息">
//	    <transition g="" name="transition4" offset="0,0" to="decision1"/>
//    </custom>
//   <decision displayName="decision1"  layout="384,118,-1,-1" name="decision1">
//       <transition displayName="同意" expr="#result" g="" name="agree" offset="0,0" to="end1"/>
//       <transition displayName="不同意" g="408,68;172,68" name="disagree" offset="0,0" to="end1"/>
//   </decision>
//</process>
//`,
//	}
//
//	err := e.Process().SaveProcess(s)
//	fmt.Println(err)
//	vv, _ := json.Marshal(s)
//	fmt.Println(string(vv))
//
//
//
//	// active
//	e.Process().ActiveProcess(1)

	fmt.Println(e.StartInstanceByIdAndOperator(1, "ll"))
}