// KeepMeBot - executor
// 2020-08-20 21:57
// Benny <benny.think@gmail.com>

package main

// sanitize is necessary

// currently supported command list
import (
	"bytes"
	"fmt"
	log "github.com/sirupsen/logrus"
	"os/exec"
)

func scheduler() {

	log.Infoln("Start scheduler")
	log.Infoln("Get all tasks")
	var tasks []Queue
	DB.Find(&tasks)
	log.Infof("Total tasks count: %d", len(tasks))
	for i, v := range tasks {
		log.Infof("Executing [%d/%d]: %s - %s(%d)", i+1, len(tasks), v.Command, v.UserName, v.UserID)

		var message string
		cmd := exec.Command("bash", "-c", v.Command)
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

		h := History{
			BaseModel: BaseModel{},
			UserID:    v.UserID,
			UserName:  v.UserName,
			Command:   v.Command,
			Output:    message,
		}
		DB.Create(&h)
	}
}
