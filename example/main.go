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

func (a *AA) Execute(context *smart.Context, model smart.INodeModel) error {
	context.Args["result"] = true
	fmt.Println("custom ....")
	return model.RunOutTransition(context)
}

func main() {

	orm.Debug = true

	e := smart.NewSmartEngine()

	//addAllCustomP(e)

	fmt.Println(e.StartInstanceByIdAndOperator(2, "ll"))
}

func addAllCustomP(engine smart.Engine) {
		s := &smart.Process{
			Name: "custom",
			DisplayName: "哈哈",
			Content: `<?xml version="1.0" encoding="UTF-8" standalone="no"?>
	
	<process displayName="借款申请流程" instanceUrl="/snaker/flow/all" name="borrow">
	  <start displayName="start1" layout="42,118,-1,-1" name="start1">
	      <transition g="" name="transition1" offset="0,0" to="load1"/>
	  </start>
	  <end displayName="end1" layout="479,118,-1,-1" name="end1"/>
	
	  
		<custom name="load1" clazz="load" displayName="获取用户信息">
		    <transition g="" name="transition4" offset="0,0" to="xx"/>
	   </custom>

<custom name="xx" clazz="load" displayName="xxx">
		    <transition g="" name="transition4" offset="0,0" to="end1"/>
	   </custom>
	 
	</process>
	`,
		}

		err := engine.Process().SaveProcess(s)
		fmt.Println(err)

		engine.Process().ActiveProcess(s.Id)
}