package smart

// smart engine 执行过程中，内容上下文
type Context struct {

	//
	Engine Engine

	// 上下文参数
	Args map[string]interface{}
}