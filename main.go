package main

import (
	"fmt"
	"github.com/andersfylling/disgord"
	"github.com/andersfylling/disgord/std"
	"os"
)

func main() {
    client := disgord.New(&disgord.Config{
        BotToken: os.Getenv("DISGORD_TOKEN"),
        Logger: disgord.DefaultLogger(false), // debug=false
    })
    if u, err := client.CreateBotURL(); err == nil {
    	fmt.Println(u)
	}

    defer client.StayConnectedUntilInterrupted()
    filter, _ := std.NewMsgFilter(client)

    client.On(disgord.EvtMessageCreate, filter.NotByBot, containsCodeBlock, sayHello)
}

