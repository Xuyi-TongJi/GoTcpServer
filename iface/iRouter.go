package iface

/*
	路由抽象接口
	路由里的苏局都是IRequest
*/

type IRouter interface {
	// PreHandle 在处理connection业务之前的钩子方法Hook
	PreHandle(request IRequest)
	// DoHandle 在处理connection业务的主方法Hook
	DoHandle(request IRequest)
	// PostHandle 在处理connection业务之后的方法Hook
	PostHandle(request IRequest)
}

// Handle 模版方法
func Handle(router IRouter, request IRequest) {
	router.PreHandle(request)
	router.DoHandle(request)
	router.PostHandle(request)
}
