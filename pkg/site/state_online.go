package site

var onlineState = make(chan *Context, 2048)

// OnlineStateChan 在线状态管道
func OnlineStateChan() *(chan *Context) {
	return &onlineState
}
