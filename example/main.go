package main

import "github.com/astaxie/beego/orm"
import _ "tianwei.pro/smart"
import _ "github.com/go-sql-driver/mysql"


func init() {
	orm.RegisterDataBase("default", "mysql", "root:xxx@tcp(127.0.0.1:3306)/smart?charset=utf8", 30)
	orm.RunSyncdb("default", false, true)
}

func main() {

	orm.Debug = true

}