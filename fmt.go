package main

import (
	"github.com/andersfylling/disgord"
	"strings"
)

func getMsg(evt interface{}) (msg *disgord.Message) {
	switch t := evt.(type) {
	case *disgord.MessageCreate:
		msg = t.Message
	case *disgord.MessageUpdate:
		msg = t.Message
	default:
		msg = nil
	}

	return msg
}

func containsCodeBlock(evt interface{}) interface{} {
	msg := getMsg(evt)
	if !strings.Contains(strings.ToLower(msg.Content), "```go") {
		return nil
	}

	return evt
}

func sayHello(s disgord.Session, data *disgord.MessageCreate) {
	_, err := data.Message.Reply(s, "hello")
	s.Logger().Error(err)
}