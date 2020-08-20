// KeepMeBot - handler
// 2020-08-23 17:38
// Benny <benny.think@gmail.com>

package main

import (
	"fmt"
	"gopkg.in/tucnak/telebot.v2"
	"strings"
)

func start(m *telebot.Message) {
	_ = b.Notify(m.Sender, telebot.Typing)
	_, _ = b.Send(m.Sender, "Keep me bot by Benny")
}

func add(m *telebot.Message) {
	var selector = &telebot.ReplyMarkup{}

	services, _ := getServiceArray()
	var btns []telebot.Btn
	for _, v := range services {
		btn := selector.Data(v.Name, v.Name)
		btns = append(btns, btn)
	}

	selector.Inline(
		selector.Row(btns...),
	)

	_ = b.Notify(m.Sender, telebot.Typing)
	_, _ = b.Send(m.Sender, "Select your services", selector)
}

func text(m *telebot.Message) {
	keepType, found := cache[m.Sender.ID]
	if found {
		status, cmd := addQueue(m, keepType)
		_, _ = b.Send(m.Sender, fmt.Sprintf("Status: %s\nCommands: `%s`", status, cmd),
			&telebot.SendOptions{ParseMode: telebot.ModeMarkdown},
		)
		delete(cache, m.Sender.ID)
	} else {
		_ = b.Notify(m.Sender, telebot.Typing)
		_, _ = b.Send(m.Sender, "How may I help you?")
	}

}

func on(c *telebot.Callback) {
	_ = b.Respond(c, &telebot.CallbackResponse{Text: "hhhh"})
	_, _ = b.Send(c.Sender, fmt.Sprintf("You choose %s, now tell me your address ", c.Data))
	// userid, message id
	cache[c.Sender.ID] = strings.Replace(c.Data, "\f", "", -1)

}
