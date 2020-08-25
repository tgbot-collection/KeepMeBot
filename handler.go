// KeepMeBot - handler
// 2020-08-23 17:38
// Benny <benny.think@gmail.com>

package main

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	tb "gopkg.in/tucnak/telebot.v2"
	"strings"
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
	cache[c.Sender.ID] = strings.Replace(c.Data, " ", "", -1)

}

func registerButtonNextStep(btn tb.Btn, fun func(c *tb.Callback)) {
	log.Infoln("Registering ", btn.Unique)
	b.Handle(&btn, fun)
}

func addDockerHub(m *tb.Message) string {
	data := Queue{
		UserID:      m.Sender.ID,
		UserName:    m.Sender.Username,
		Command:     m.Text,
		ServiceType: "Docker Hub",
	}
	if err := DB.Create(&data).Error; err != nil {
		return err.Error()
	} else {
		return "Your command has been add to Queue"
	}
}
func onText(m *tb.Message) {

	if value, found := cache[m.Sender.ID]; found {
		var message string
		switch value {
		case "GitHub":
			message = addDockerHub(m)
		case "DockerHub":
			message = addDockerHub(m)
		default:
			message = "Not register action"
		}
		delete(cache, m.Sender.ID)
		_, _ = b.Send(m.Sender, message)

	} else {
		_, _ = b.Send(m.Sender, "hello there")

	}
}
