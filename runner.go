// KeepMeBot - executor
// 2020-08-20 21:57
// Benny <benny.think@gmail.com>

package main

// sanitize is necessary

// currently supported command list
import (
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

		out, err := exec.Command("bash", "-c", v.Command).Output()
		if err != nil {
			message := fmt.Sprintf("Command failed %s", err)
			log.Warningln(message)
			out = []byte(message)
		} else {
			log.Infof("%s", out)
		}
		h := History{
			BaseModel: BaseModel{},
			UserID:    v.UserID,
			UserName:  v.UserName,
			Command:   v.Command,
			Output:    string(out),
		}
		DB.Create(&h)
	}
}
