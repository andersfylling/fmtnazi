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
	code, err := gofmt([]byte(data.Message.Content))
	if err != nil {
		_, err := data.Message.Reply(s, err.Error())
		s.Logger().Error(err)
		return
	}

	reply := string(code)
	_, err = data.Message.Reply(s, reply)
	s.Logger().Error(err)
}

func gofmt(content []byte) ([]byte, error) {
	start := strings.Index(string(content), "```go") + len("```go")
	code := content[start:]
	end := strings.Index(string(code), "```")
	code = code[:end]
	codeStr := string(code)

	if !strings.Contains(codeStr, "func main") {
		return nil, errors.New("missing func main")
	}
	if strings.Count(codeStr, "\n") < 3 {
		return nil, errors.New("there were less than 3 new lines")
	}

	// go fmt
	formatted, err := format.Source(code)
	if err != nil {
		return nil, err
	}

	// replace \t with spaces
	const tabInSpaces = "    "
	const tab = "\t"
	fStr := strings.Replace(string(formatted), tab, tabInSpaces, -1)
	formatted = []byte(fStr)

	// wrap in code block
	formatted = append([]byte("```go"), formatted...)
	formatted = append(formatted, []byte("```")...)

	return formatted, nil
}
