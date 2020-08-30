// KeepMeBot - database_test
// 2020-08-28 20:11
// Benny <benny.think@gmail.com>

package main

import (
	"crypto/md5"
	"fmt"
	"github.com/stretchr/testify/assert"
	"gopkg.in/tucnak/telebot.v2"
	"os"
	"testing"
)
import "gopkg.in/alessio/shellescape.v1"

func TestMain(m *testing.M) {
	setup()
	code := m.Run()
	teardown()
	// or rm test.db in makefile.
	os.Exit(code)
}

func setup() {
	// Do something here.
	_ = os.Setenv("test", "true")
	deferInit()
	fmt.Printf("\033[1;36m%s\033[0m", "> Setup completed\n")
}

func teardown() {
	// Do something here.
	_ = os.Remove("test.db")
	fmt.Printf("\033[1;36m%s\033[0m", "> Teardown completed")
	fmt.Printf("\n")
}

func TestDeferInit(t *testing.T) {
	count := 0
	DB.Model(&Service{}).Count(&count)
	assert.NotEqual(t, count, 0)
}

func TestGetServiceArray(t *testing.T) {
	arr, err := getServiceArray()
	assert.Nil(t, err)
	for i, ser := range arr {
		assert.Equal(t, ser.Name, supportedService[i].Name)
	}
}

func TestGetServiceMap(t *testing.T) {
	smap := getServiceMap()
	var names []string
	for _, item := range supportedService {
		names = append(names, item.Name)
	}
	for i := range smap {
		assert.True(t, inArray(i, names))
	}
}

func TestAddQueue(t *testing.T) {
	//m telebot.Message, serviceName string, inputs ...interface{}
	var docker = telebot.Message{
		Sender: &telebot.User{
			ID:       1234,
			Username: "BennyThink",
		},
		Text: shellescape.Quote("alpine"),
	}
	var github = telebot.Message{
		Sender: &telebot.User{
			ID:       1234,
			Username: "BennyThink",
		},
		Text: "https://github.com/BennyThink/KeepMeBot/",
	}
	var http = telebot.Message{
		Sender: &telebot.User{
			ID:       1234,
			Username: "BennyThink",
		},
		Text: "http://z.cn",
	}
	a := getServiceMap()

	var text, message string
	text = shellescape.Quote(docker.Text)
	message = addQueue(docker, "Docker Hub", text, text)
	assert.Contains(t, message, fmt.Sprintf(a["Docker Hub"].Template, text, text))

	text = shellescape.Quote(docker.Text)
	dest := fmt.Sprintf("%x", md5.Sum([]byte(text)))
	message = addQueue(github, "GitHub", text+" "+dest, dest)
	assert.Contains(t, message, fmt.Sprintf(a["GitHub"].Template, text+" "+dest, dest))

	text = shellescape.Quote(http.Text)
	message = addQueue(http, "get", text)
	assert.Contains(t, message, fmt.Sprintf(a["get"].Template, text))

	// test add too many for another user
	var another = telebot.Message{
		Sender: &telebot.User{
			ID:       9876,
			Username: "someone else",
		},
		Text: shellescape.Quote("nginx"),
	}
	for i := 0; i <= 5; i++ {
		addQueue(another, "Docker Hub", text, text)
	}
	errMsg := addQueue(another, "Docker Hub", text, text)
	assert.Contains(t, errMsg, "Your limit")
}

func TestGetQueueList(t *testing.T) {
	data := getQueueList(1234)
	assert.Equal(t, 3, len(data))
	data = getQueueList(9876)
	assert.Equal(t, 5, len(data))
}

func TestDeleteQueue(t *testing.T) {
	deleteQueue("6")
	data := getQueueList(9876)
	assert.Equal(t, 4, len(data))
}

func TestSession(t *testing.T) {
	setSession(1234, "test-next")

	next := getSession(1234)
	assert.Equal(t, next, "test-next")

	next = getSession(1234444)
	assert.Empty(t, next)

	deleteSession(1234)
	next = getSession(1234)
	assert.Empty(t, next)
}

func TestHistoryRecorder(t *testing.T) {
	var v = Queue{
		UserID:    1234,
		UserName:  "BennyThink",
		Command:   "history command",
		ServiceID: 1,
	}
	historyRecorder(v, "message")
}

func TestGetHistory(t *testing.T) {
	h := getHistory(1234)
	assert.Equal(t, "history command", h[0].Command)
}

func inArray(element string, arr []string) bool {
	for _, v := range arr {
		if v == element {
			return true
		}
	}
	return false
}
