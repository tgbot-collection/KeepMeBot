// KeepMeBot - handler
// 2020-08-23 17:38
// Benny <benny.think@gmail.com>

package main

import (
	"crypto/md5"
	"fmt"

	log "github.com/sirupsen/logrus"
	"gopkg.in/alessio/shellescape.v1"

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
		btn := selector.Data(v.Name, fmt.Sprintf("AddServiceButton%d%d", v.ID, m.Sender.ID), v.Name)
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

func removeServiceButton(c *tb.Callback) {
	_ = b.Respond(c, &tb.CallbackResponse{Text: "Ok"})
	deleteQueue(c.Data)
	_, _ = b.Send(c.Sender, fmt.Sprintf("Id %s has been deleted ", c.Data))
}

func registerButtonNextStep(btn tb.Btn, fun func(c *tb.Callback)) {
	log.Infoln("Registering ", btn.Unique)
	b.Handle(&btn, fun)
}

func onCallback(c *tb.Callback) {
	_ = b.Respond(c, &tb.CallbackResponse{Text: "You seem to delete an outdated button"})
	_, _ = b.Send(c.Sender, "You seem to delete an outdated button")
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
	text := shellescape.Quote(m.Text)
	message = addQueue(m.Sender.ID, m.Sender.Username, "Docker Hub", text, text)
	deleteSession(m.Sender.ID)
	return
}

func github(m *tb.Message) (message string) {
	text := shellescape.Quote(m.Text)
	dest := fmt.Sprintf("%x", md5.Sum([]byte(text)))
	message = addQueue(m.Sender.ID, m.Sender.Username, "GitHub", text+" "+dest, dest)
	deleteSession(m.Sender.ID)
	return
}

func list(m *tb.Message) {
	queue := getQueueList(m.Sender.ID)

	if len(queue) == 0 {
		_, _ = b.Send(m.Sender, "you dont seem to have any services")
		return
	}
	var inlineKeys [][]tb.InlineButton
	for i, v := range queue {
		var temp = tb.InlineButton{
			Unique: fmt.Sprintf("DeleteServiceButton%d%d", v.ID, m.Sender.ID),
			Text:   fmt.Sprintf("%d. %s", i+1, v.Command),
			Data:   fmt.Sprintf("%d", v.ID),
		}
		var b = tb.Btn{
			Unique: temp.Unique,
			Text:   temp.Text,
			Data:   temp.Data,
		}
		registerButtonNextStep(b, removeServiceButton)
		inlineKeys = append(inlineKeys, []tb.InlineButton{temp})
	}

	_ = b.Notify(m.Sender, tb.Typing)
	_, _ = b.Send(m.Sender, "Select to delete service", &tb.ReplyMarkup{
		InlineKeyboard: inlineKeys,
	})

}
