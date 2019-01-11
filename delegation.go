package smart

//
type Delegation interface {
	//
	Execute(context *Context, model INodeModel) error
}