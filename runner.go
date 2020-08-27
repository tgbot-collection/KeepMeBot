// KeepMeBot - executor
// 2020-08-20 21:57
// Benny <benny.think@gmail.com>

package main

// sanitize is necessary

// currently supported command list
import (
	"bytes"
	"fmt"
	"net/http"
	"os/exec"

	log "github.com/sirupsen/logrus"
)

func scheduler() {

	log.Infoln("Start scheduler")
	log.Infoln("Get all tasks")
	var tasks []Queue
	DB.Find(&tasks)
	log.Infof("Total tasks count: %d", len(tasks))
	for i, v := range tasks {
		log.Infof("[%s]Executing [%d/%d]: %s - %s(%d)", v.Service.ServiceType,
			i+1, len(tasks), v.Command, v.UserName, v.UserID)
		var message string

		switch v.Service.ServiceType {
		case "internal":
			internalExecutor(v)
		case "external":
			message = externalExecutor(v)

		}

		historyRecorder(v, message)
	}
}

func externalExecutor(q Queue) string {
	var message string
	cmd := exec.Command("bash", "-c", q.Command)
	var out bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr
	err := cmd.Run()
	if err != nil {
		message = fmt.Sprintf("%s %s", err, stderr.String())
		log.Warningln(message)
	} else {
		message = out.String()
		log.Infof(message)
	}
	return message
}

func internalExecutor(q Queue) (s string) {
	switch q.Service.Name {
	case "get":
		s = get(q.Parameter)
	}
	return
}

func get(url string) (msg string) {
	resp, err := http.Get(url)
	if err != nil {
		msg = err.Error()
	} else {
		msg = fmt.Sprintf("%s", resp.StatusCode)
		_ = resp.Body.Close()
	}
	return
}
