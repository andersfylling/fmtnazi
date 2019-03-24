package main

import (
	"bytes"
	"errors"
	"github.com/andersfylling/disgord"
	"go/format"
	"strings"
)

const prefixCodeBlock = "```go"
const suffixCodeBlock = "```"

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
	if !strings.Contains(strings.ToLower(msg.Content), prefixCodeBlock) {
		return nil
	}

	// if there is a start, make sure there is a end
	start := strings.Index(strings.ToLower(msg.Content), prefixCodeBlock)
	end := strings.Index(strings.ToLower(msg.Content[start+len(prefixCodeBlock):]), suffixCodeBlock)
	if end == -1 {
		return nil
	}

	return evt
}

func formatMessage(s disgord.Session, data *disgord.MessageCreate) {
	replyBytes, err := craftReply(bytes.Runes([]byte(data.Message.Content)))
	if err != nil {
		s.Logger().Error(err)
		return
	}

	reply := "\n#> Written by " + data.Message.Author.Mention() + "\n\n" + string(replyBytes)
	_, err = data.Message.Reply(s, reply)
	s.Logger().Error(err)
}

func craftReply(content []rune) (reply []rune, err error) {
	tmp := content
	for {
		start := strings.Index(string(tmp), prefixCodeBlock)
		if start == -1 {
			// all code blocks have been handled
			break
		}
		end := strings.Index(string(tmp[start+len(prefixCodeBlock):]), suffixCodeBlock)

		reply = append(reply, tmp[:start]...)
		code, err := gofmt(tmp[start : start+len(prefixCodeBlock)+end+len(suffixCodeBlock)])
		if err != nil {
			return nil, err
		}
		reply = append(reply, code...)
		tmp = tmp[start+len(prefixCodeBlock)+end+len(suffixCodeBlock):]

	}

	reply = append(reply, tmp...)
	return reply, nil
}

// gofmt takes a code block, with or without discord wraps, and returns
// correctly formatted go code using spaces instead of tabs.
func gofmt(content []rune) ([]rune, error) {
	var code []rune

	// unwrap
	start := strings.Index(string(content), prefixCodeBlock) + len(prefixCodeBlock)
	code = content[start:]
	if start > 0 {
		end := strings.Index(string(code), suffixCodeBlock)
		code = code[:end]
	}

	// control checks
	codeStr := string(code)
	if !strings.Contains(codeStr, "func main") {
		return nil, errors.New("missing func main")
	}
	if strings.Count(codeStr, "\n") < 3 {
		return nil, errors.New("there were less than 3 new lines")
	}

	// go fmt
	formatted, err := format.Source([]byte(codeStr))
	if err != nil {
		return nil, err
	}

	// replace \t with spaces
	const tabInSpaces = "    "
	const tab = "\t"
	fStr := strings.Replace(string(formatted), tab, tabInSpaces, -1)
	formatted = []byte(fStr)

	// wrap in code block
	if start > 0 {
		formatted = append([]byte(prefixCodeBlock), formatted...)
		formatted = append(formatted, []byte(suffixCodeBlock)...)
	}

	return bytes.Runes(formatted), nil
}
