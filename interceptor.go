package smart

// 过滤器
type Interceptor interface {

	//
	Intercept(context *Context) error
}