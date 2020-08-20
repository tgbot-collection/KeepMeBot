package main

import (
	"log"
	"os"
	"time"

	tb "gopkg.in/tucnak/telebot.v2"
)

var token = os.Getenv("token")
var b, err = tb.NewBot(tb.Settings{
	Token:  token,
	Poller: &tb.LongPoller{Timeout: 10 * time.Second},
})

func main() {

	if err != nil {
		log.Fatal(err)
		return
	}

	b.Handle("/hello", funcName)
	b.Start()
}

func funcName(m *tb.Message) {
	_, _ = b.Send(m.Sender, "Hello World!")
}
