package ziface

/*
路由接口：
路由里的数据都是IRequest
*/
type IRouter interface {
	PreHandle(request IRequest)  // 在处理conn业务之前的方法hook
	Handle(request IRequest)     // 在处理conn业务的主方法hook
	PostHandle(request IRequest) // 在处理conn业务之后的方法hook
}
