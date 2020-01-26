package wslistner

import (
	"fmt"
	"sync"

	"github.com/gorilla/websocket"
)

type Conn interface {
	Run(readCh chan []byte, closeCh chan bool)
	Write([]byte)
	Close()
}

type conn struct {
	ws      *websocket.Conn
	wg      *sync.WaitGroup
	writeCh chan []byte
}

func NewConn(ws *websocket.Conn) Conn {
	return &conn{
		ws: ws,
	}
}

func (c *conn) Run(readCh chan []byte, closeCh chan bool) {
	c.wg = &sync.WaitGroup{}
	c.writeCh = make(chan []byte)

	errCh := make(chan error)
	c.wg.Add(1)
	go c.waitRead(readCh, errCh)

	c.wg.Add(1)
	go c.waitWrite()

	for {
		select {
		case <-errCh:
			close(c.writeCh)
			c.wg.Wait()

			close(closeCh)
			return
		}
	}
}

func (c *conn) Write(bytes []byte) {
	c.writeCh <- bytes
}

func (c *conn) Close() {
	if err := c.ws.Close(); err != nil {
		fmt.Println(err.Error())
	}
}

func (c *conn) waitWrite() {
	defer c.wg.Done()

	for bytes := range c.writeCh {
		if err := c.ws.WriteMessage(websocket.TextMessage, bytes); err != nil {
			fmt.Println(err.Error())
			break
		}
	}
	c.Close()
}

func (c *conn) waitRead(readCh chan []byte, errCh chan error) {
	defer c.wg.Done()

	for {
		_, readBytes, err := c.ws.ReadMessage()
		if err != nil {
			errCh <- err
			break
		}
		readCh <- readBytes
	}
	c.Close()
}