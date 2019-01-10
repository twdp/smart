package main

import (
	"encoding/json"
	"fmt"
	"github.com/astaxie/beego/orm"
	"tianwei.pro/smart"
)
import _ "tianwei.pro/smart"
import _ "github.com/go-sql-driver/mysql"


func init() {
	//orm.DefaultTimeLoc = time.Local

	orm.RegisterDataBase("default", "mysql", "root:xxx@tcp(127.0.0.1:3306)/smart?charset=utf8&loc=Asia%2FShanghai", 30)
	orm.RunSyncdb("default", false, true)
}

func main() {

	orm.Debug = true

	e := smart.NewSmartEngine()

	s := &smart.Process{
		Name: "xx",
		DisplayName: "ccc",
		Content: `<?xml version="1.0" encoding="UTF-8" standalone="no"?>
	
			<process displayName="借款申请流程" instanceUrl="/snaker/flow/all" name="borrow">
			   <start displayName="start1" layout="42,118,-1,-1" name="start1">
			       <transition g="" displayName="xx" name="transition1" offset="0,0" to="apply"/>
			   </start>
			   <end displayName="end1" layout="479,118,-1,-1" name="end1"/>
	
			   <task assignee="apply.operator" autoExecute="Y" displayName="借款申请" form="/flow/borrow/apply" layout="126,116,-1,-1" name="apply" performType="ANY" taskType="Major">
			       <transition g="" displayName="xx" name="transition2" offset="0,0" to="approval"/>
			   </task>
			   <task assignee="approval.operator" autoExecute="Y" displayName="审批" form="/snaker/flow/approval" layout="252,116,-1,-1" name="approval" performType="ANY" taskType="Major">
			    <transition g="" displayName="xx" name="transition3" offset="0,0" to="decision1"/>
			   </task>
	
			</process>`,
	}

	err := e.Process().SaveProcess(s)
	fmt.Println(err)
	vv, _ := json.Marshal(s)
	fmt.Println(string(vv))



	// active
	e.Process().ActiveProcess(1)


}