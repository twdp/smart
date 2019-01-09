package smart

// smart engine 执行过程中，内容上下文
type Context struct {

	// 当前引擎执行的实例
	Instance *Instance

	// 执行引擎
	Engine Engine

	// 上下文参数
	Args map[string]interface{}

	/**
	 * 是否已合并
	 * 针对join节点的处理
	 */
	 IsMerged bool
}