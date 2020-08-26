// KeepMeBot - handler
// 2020-08-23 17:38
// Benny <benny.think@gmail.com>

package main

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	tb "gopkg.in/tucnak/telebot.v2"
)

func start(m *tb.Message) {
	_ = b.Notify(m.Sender, tb.Typing)
	_, _ = b.Send(m.Sender, "Keep me bot by Benny")
}

func add(m *tb.Message) {
	var selector = &tb.ReplyMarkup{}

	services, _ := getServiceArray()
	var btns []tb.Btn
	for _, v := range services {
		btn := selector.Data(v.Name, fmt.Sprintf("AddServiceButton%d", v.ID), v.Name)
		registerButtonNextStep(btn, addServiceButton)
		btns = append(btns, btn)
	}

	selector.Inline(
		selector.Row(btns...),
	)

	_ = b.Notify(m.Sender, tb.Typing)
	_, _ = b.Send(m.Sender, "Select your services", selector)
}

func addServiceButton(c *tb.Callback) {
	_ = b.Respond(c, &tb.CallbackResponse{Text: "Ok"})
	_, _ = b.Send(c.Sender, fmt.Sprintf("You choose %s, now tell me your address ", c.Data))
	setSession(c.Sender.ID, c.Data)
}

func registerButtonNextStep(btn tb.Btn, fun func(c *tb.Callback)) {
	log.Infoln("Registering ", btn.Unique)
	b.Handle(&btn, fun)
}

func onText(m *tb.Message) {
	current := getSession(m.Sender.ID)
	var message string
	switch current {
	case "Docker Hub":
		message = dockerhub(m)
	case "GitHub":
		message = github(m)
	default:
		message = "hello"
	}
	_, _ = b.Send(m.Sender, message, &tb.SendOptions{
		ParseMode: tb.ModeMarkdown,
	})

}

func dockerhub(m *tb.Message) (message string) {
	message = addQueue(m.Sender.ID, m.Sender.Username, m.Text, "Docker Hub")
	deleteSession(m.Sender.ID)
	return
}

func github(m *tb.Message) (message string) {
	message = addQueue(m.Sender.ID, m.Sender.Username, m.Text, "GitHub")
	deleteSession(m.Sender.ID)
	return
}
