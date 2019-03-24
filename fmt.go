package main

import (
	"errors"
	"github.com/andersfylling/disgord"
	"go/format"
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

func formatCode(s disgord.Session, data *disgord.MessageCreate) {
	code, err := gofmt(data.Message.Content)
	if err != nil {
		_, err := data.Message.Reply(s, err.Error())
		s.Logger().Error(err)
		return
	}

	reply := "```go\n" + code + "\n```"
	_, err = data.Message.Reply(s, reply)
	s.Logger().Error(err)
}

func gofmt(str string) (string, error) {
	start := strings.Index(str, "```go") + len("```go")
	code := str[start:]
	end := strings.Index(code, "```")
	code = code[:end]

	if !strings.Contains(code, "func main") {
		return "", errors.New("missing func main")
	}
	if strings.Count(code, "\n") < 3 {
		return "", errors.New("there were less than 3 new lines")
	}

	formatted, err := format.Source([]byte(code))
	if err != nil {
		return "", err
	}

	return string(formatted), nil
}
