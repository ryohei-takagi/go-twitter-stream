package app

import (
	"github.com/flatnyat/go-twitter-stream/pkg/wslistner"
)

type Router interface {
	Run()
}

type router struct {
	port int
	twitter *TwitterClient
}

func NewRouter(port int, twitter *TwitterClient) Router {
	return &router{
		port,
		twitter,
	}
}

func (r *router) Run() {
	wsListener := wslistner.NewListener(r.port)
	wsListener.RegisterAcceptHandler(r.OnAccept)
	wsListener.RegisterCloseHandler(r.OnClose)
	wsListener.Run()
}

func (r *router) OnAccept(c wslistner.Conn) {
	u := NewUser(c, r.twitter)
	u.Run()
}

func (r *router) OnClose(c wslistner.Conn) {
}