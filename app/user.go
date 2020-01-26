package app

import (
	"encoding/json"
	"fmt"
	"github.com/flatnyat/go-twitter-stream/pkg/wslistner"
)

type User interface {
	Run()
	Write(body interface{})
	Payload() *Payload
}

type user struct {
	conn   wslistner.Conn
	readCh chan []byte
	twitter *TwitterClient
	payload *Payload
}

func NewUser(c wslistner.Conn, t *TwitterClient) User {
	return &user{
		conn:   c,
		readCh: make(chan []byte),
		twitter: t,
	}
}

func (u *user) Run() {
	readCh := make(chan []byte)
	closeCh := make(chan bool)

	go u.conn.Run(readCh, closeCh)

	// Waiting to process message...
	for {
		select {
		case bytes := <-readCh:
			u.doHandler(bytes)

		case <-closeCh:
			u.doClose()
			return

		default:
		}
	}
}

func (u *user) Write(body interface{}) {
	payload := &Payload{
		Body: body,
	}
	bytes, err := json.Marshal(payload)
	if err != nil {
		fmt.Println(err.Error())
		u.conn.Close()
	}
	u.conn.Write(bytes)
}

func (u *user) Payload() *Payload {
	return u.payload
}

// Received message from Client
func (u *user) doHandler(bytes []byte) {
	payload := &Payload{}
	if err := json.Unmarshal(bytes, payload); err != nil {
		fmt.Println(err.Error())
		u.conn.Close()
	}
	u.payload = payload

	s := *u.twitter
	if err := s.Connect(u.payload, u.handleMessage); err != nil {
		fmt.Println(err.Error())
		u.conn.Close()
	}
}

func (u *user) doClose() {
	u.payload.Body = ""
	t := *u.twitter
	if err := t.Shutdown(u.payload, u.handleMessage); err != nil {
		fmt.Println(err.Error())
	}
	u.conn.Close()
}

func (u *user) handleMessage(result interface{}) {
	fmt.Println(result)
	// u.Write(result)
}
