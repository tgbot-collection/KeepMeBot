// KeepMeBot - handler
// 2020-08-23 17:38
// Benny <benny.think@gmail.com>

package main

import (
	"crypto/md5"
	"fmt"

	log "github.com/sirupsen/logrus"
	"github.com/tgbot-collection/tgbot_ping"
	"gopkg.in/alessio/shellescape.v1"
	tb "gopkg.in/tucnak/telebot.v2"
)

var messageMap = map[string]string{
	"Docker Hub": "now tell me your repository, example `bennythink/keepmebot`",
	"GitHub":     "now tell me your repository, example `https://github.com/BennyThink/KeepMeBot/`",
	"get":        "now tell me your url, example `https://github.com/BennyThink/KeepMeBot/`",
}

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
	_, _ = b.Send(c.Sender,
		fmt.Sprintf("You choose %s, %s ", c.Data, messageMap[c.Data]),
		&tb.SendOptions{
			ParseMode: tb.ModeMarkdownV2,
		})
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
	case "get":
		message = getFunc(m)
	default:
		message = "this session hasn't registered."
	}
	_, _ = b.Send(m.Sender, message, &tb.SendOptions{
		ParseMode: tb.ModeMarkdownV2,
	})

}

func dockerhub(m *tb.Message) (message string) {
	text := shellescape.Quote(m.Text)
	message = addQueue(*m, "Docker Hub", text, text)
	deleteSession(m.Sender.ID)
	return
}

func github(m *tb.Message) (message string) {
	text := shellescape.Quote(m.Text)
	dest := fmt.Sprintf("%x", md5.Sum([]byte(text)))
	message = addQueue(*m, "GitHub", text+" "+dest, dest)
	deleteSession(m.Sender.ID)
	return
}

func getFunc(m *tb.Message) (message string) {
	text := shellescape.Quote(m.Text)
	message = addQueue(*m, "get", text)
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

func history(m *tb.Message) {
	h := getNewestHistory(m.Sender.ID)
	_ = b.Notify(m.Sender, tb.Typing)
	if len(h) == 0 {
		_, _ = b.Send(m.Sender, "A brave new world!")
		return
	}

	for i, v := range h {
		message := fmt.Sprintf("%d. `%s`\n %s | %s", i+1, v.Command,
			v.CreatedAt.Format("2006-01-02 15:04:05"), v.Output)

		_, _ = b.Send(m.Sender, message, &tb.SendOptions{
			ParseMode: tb.ModeMarkdown,
		})
	}

}

func ping(m *tb.Message) {
	_ = b.Notify(m.Sender, tb.Typing)
	bot := tgbot_ping.GetRuntime("botsrunner_keepme-runner_1", "KeepMeBot", "html")
	docker := tgbot_ping.GetRuntime("botsrunner_keepme-docker_1", "KeepMeDocker", "html")
	_, _ = b.Send(m.Chat, bot+"\n\n"+docker, &tb.SendOptions{ParseMode: tb.ModeHTML})
}

func edited(m *tb.Message) {
	_ = b.Notify(m.Sender, tb.RecordingVNote)
	var message = `You can edit your message, you can rewrite history in Git - don't let them do it for real!
Remember, remember!
`
	_, _ = b.Send(m.Sender, message)
}
