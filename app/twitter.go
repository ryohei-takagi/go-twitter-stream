package app

import (
	"github.com/dghubble/go-twitter/twitter"
	"github.com/dghubble/oauth1"
	"time"
)

type TwitterClient interface {
	Connect(*Payload, AcceptStreamFunc) error
	Reboot(*Payload, AcceptStreamFunc) error
	Shutdown(*Payload, AcceptStreamFunc) error
}

type twitterClient struct {
	client *twitter.Client
	stream *StreamClient
}

func NewTwitterClient() TwitterClient {
	config := oauth1.NewConfig("", "")
	token := oauth1.NewToken("", "")

	httpClient := config.Client(oauth1.NoContext, token)

	client := twitter.NewClient(httpClient)

	return &twitterClient{
		client: client,
		stream: nil,
	}
}

func (t *twitterClient) Connect(payload *Payload, fn AcceptStreamFunc) error {
	// 1人目の接続
	if t.stream == nil {
		keyword := Keyword{
			ID:    payload.ID,
			Value: payload.Body.(string),
		}
		return t.streaming([]Keyword{keyword}, fn)
	}

	// 2人目以降の接続
	return t.Reboot(payload, fn)
}

func (t *twitterClient) Reboot(payload *Payload, fn AcceptStreamFunc) error {
	stream := *t.stream

	// TwitterAPIの制限を回避するため5秒程度止める
	stream.Stop()
	time.Sleep(5 * time.Second)

	// 自身以外のキーワードを再セット
	var keywords []Keyword
	for _, v := range stream.Keywords() {
		if v.ID != payload.ID {
			keywords = append(keywords, v)
		}
	}

	// 自身のキーワードをセット
	if payload.Body != "" {
		keyword := Keyword{
			ID:    payload.ID,
			Value: payload.Body.(string),
		}
		keywords = append(keywords, keyword)
	}

	return t.streaming(keywords, fn)
}

func (t *twitterClient) Shutdown(payload *Payload, fn AcceptStreamFunc) error {
	if t.stream == nil {
		return nil
	}

	return t.Reboot(payload, fn)
}

func (t *twitterClient) streaming(keywords []Keyword, fn AcceptStreamFunc) error {
	stream := NewStreamClient(t, keywords)
	t.stream = &stream

	go stream.Subscribe(fn)

	return nil
}