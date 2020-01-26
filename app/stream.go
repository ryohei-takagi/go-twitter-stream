package app

import (
	"fmt"
	"github.com/dghubble/go-twitter/twitter"
)

type AcceptStreamFunc func(interface{})

type StreamClient interface {
	Keywords() []Keyword
	Subscribe(AcceptStreamFunc)
	Stop()
}

type streamClient struct {
	client *twitter.Stream
	keywords []Keyword
}

type Keyword struct {
	ID int
	Value string
}

func NewStreamClient(t *twitterClient, keywords []Keyword) StreamClient {
	var values []string
	for _, k := range keywords {
		values = append(values, k.Value)
	}

	params := &twitter.StreamFilterParams{
		Track: values,
		StallWarnings: twitter.Bool(true),
	}

	client, err := t.client.Streams.Filter(params)
	if err != nil {
		fmt.Println(err.Error())
		return nil
	}
fmt.Println(keywords)
	return &streamClient{
		client:   client,
		keywords: keywords,
	}
}

func (s *streamClient) Keywords() []Keyword {
	return s.keywords
}

func (s *streamClient) Subscribe(fn AcceptStreamFunc) {
	for message := range s.client.Messages {
		fn(message)
	}
}

func (s *streamClient) Stop() {
	s.client.Stop()
}