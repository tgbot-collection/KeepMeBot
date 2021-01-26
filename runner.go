// KeepMeBot - executor
// 2020-08-20 21:57
// Benny <benny.think@gmail.com>

package main

// sanitize is necessary

// currently supported command list
import (
	"bytes"
	"fmt"
	"math"
	"net/http"
	"os/exec"
	"time"

	log "github.com/sirupsen/logrus"
)

func scheduler() {
	log.Debugln("Start scheduler")
	var tasks []Queue
	var executeTasks []Queue
	DB.Find(&tasks)
	log.Debugf("Total prepare tasks count: %d", len(tasks))

	for _, item := range tasks {
		var h History
		var s Service
		DB.Where("service_id=?", item.ServiceID).
			Order("created_at desc").First(&h)
		DB.Find(&h).Related(&s)
		if time.Now().Sub(h.CreatedAt) > time.Duration(s.Interval*math.Pow10(9)) ||
			h.UserID == 0 {
			log.Debugf("Add   %s(%v):%s to queue...", item.UserName, item.UserID, item.Command)
			executeTasks = append(executeTasks, item)
		}
	}
	log.Debugf("Total execute tasks count: %d", len(executeTasks))

	for i, v := range executeTasks {
		log.Infof("Executing [%d/%d]", i+1, len(tasks))
		var message string
		var s Service
		DB.Model(&v).Related(&s)
		switch s.ServiceType {
		case "internal":
			message = internalExecutor(v)
		case "external":
			message = externalExecutor(v)
		default:
			log.Warningln("404")
		}

		historyRecorder(v, message)
	}
}

func externalExecutor(q Queue) string {
	log.Infof("External Job for %s(%v):%s", q.UserName, q.UserID, q.Command)

	var message string
	cmd := exec.Command("sh", "-c", q.Command)
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
		log.Debugf(message)
	}
	return message
}

func internalExecutor(q Queue) (s string) {
	log.Infof("Internal Job for %s(%v):%s", q.UserName, q.UserID, q.Command)
	var ser Service
	DB.Model(&q).Related(&ser)
	switch ser.Name {
	case "get":
		s = get(q.Parameter)
	default:
		log.Warningln("Internal job not found.")
	}
	return
}

func get(url string) (msg string) {
	resp, err := http.Get(url)
	if err != nil {
		msg = err.Error()
	} else {
		msg = fmt.Sprintf("http status code is %d", resp.StatusCode)
		_ = resp.Body.Close()
	}
	return
}
