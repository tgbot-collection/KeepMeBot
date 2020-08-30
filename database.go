// KeepMeBot - database
// 2020-08-20 22:01
// Benny <benny.think@gmail.com>

package main

import (
	"fmt"
	"github.com/jinzhu/gorm"
	log "github.com/sirupsen/logrus"
	"gopkg.in/tucnak/telebot.v2"
	"os"
	"time"
)

// SQLite

// all the support keep alive service are bundle with database.

var DB *gorm.DB

type BaseModel struct {
	gorm.Model
}

type Service struct {
	BaseModel
	Name        string `gorm:"unique"`
	Max         int
	ServiceType string `gorm:"default: external"`
	Template    string
	Interval    float64 `gorm:"default: 86400"` // seconds
}

type Queue struct {
	BaseModel
	UserID    int
	UserName  string
	Parameter string
	Command   string
	Service   Service
	ServiceID int
}

type History struct {
	BaseModel
	UserID    int
	UserName  string
	Command   string
	Output    string
	Service   Service
	ServiceID int
}

type Session struct {
	BaseModel
	UserID int `gorm:"unique"`
	Next   string
}

var supportedService []Service

func deferInit() {
	supportedService = []Service{
		{
			Name:        "Docker Hub",
			Max:         5,
			ServiceType: "external",
			Template:    "docker pull %s && docker rmi %s",
		},
		{
			Name:        "GitHub",
			Max:         3,
			ServiceType: "external",
			Template:    "git clone %s && rm -rf %s",
		},
		{
			Name:        "get",
			Max:         10,
			ServiceType: "internal",
			Template:    "get %s",
			Interval:    time.Second.Seconds() * 60,
		},
	}
	var err error
	var dbFile string

	switch os.Getenv("test") {
	case "true":
		dbFile = "test.db"
	default:
		dbFile = "keep.db"
	}
	log.Infof("Using %s as database", dbFile)
	DB, err = gorm.Open("sqlite3", dbFile)
	if err != nil {
		panic("failed to connect database")
	}

	// Migrate the schema
	DB.AutoMigrate(&Service{}, &Queue{}, &History{}, &Session{})

	for _, v := range supportedService {
		DB.FirstOrCreate(&v, v)
	}

}

func getServiceArray() (supportedService []Service, err error) {
	if err = DB.Find(&supportedService).Error; err != nil {
		return
	}
	return
}

func getQueueList(userid int) (q []Queue) {
	DB.Where("user_id=?", userid).Find(&q)
	return
}

func deleteQueue(qid string) {
	DB.Where("id=?", qid).Delete(&Queue{})
}

func getServiceMap() map[string]Service {
	var a = make(map[string]Service)
	arr, _ := getServiceArray()

	for _, v := range arr {
		a[v.Name] = v
	}
	return a
}

func addQueue(m telebot.Message, serviceName string, inputs ...interface{}) (message string) {
	// user id, commands
	s := getServiceMap()
	service := s[serviceName]

	realCommand := fmt.Sprintf(service.Template, inputs...)
	// max than 5
	count := 0
	DB.Model(&Queue{}).Where("user_id=?", m.Sender.ID).Count(&count)
	log.Infof("Check for %v, current %d", m.Sender.ID, count)
	if count > service.Max-1 {
		message = fmt.Sprintf("Your limit is %d, you are using %d", service.Max, count)
	} else {
		d := Queue{
			UserID:    m.Sender.ID,
			UserName:  m.Sender.Username,
			Parameter: m.Text,
			Command:   realCommand,
			Service:   service,
		}

		DB.Create(&d)
		message = fmt.Sprintf("%s Your command is `%s`", "Success\\!", realCommand)
	}
	return
}

func setSession(id int, next string) {
	// create or update
	session := Session{
		UserID: id,
		Next:   next,
	}
	DB.Save(&session)

}

func getSession(id int) string {
	session := Session{}
	DB.Where("user_id=?", id).First(&session)
	return session.Next
}

func deleteSession(id int) {
	session := Session{
		UserID: id,
	}
	DB.Unscoped().Delete(&session)
}

func historyRecorder(v Queue, message string) {
	h := History{
		BaseModel: BaseModel{},
		UserID:    v.UserID,
		UserName:  v.UserName,
		Command:   v.Command,
		Output:    message,
		ServiceID: v.ServiceID,
	}
	DB.Create(&h)
}

func getHistory(userId int) (h []History) {
	//select *, max(created_at)
	//from histories
	//where user_id = 260260121
	//group by service_id
	DB.Where("user_id=?", userId).Group("service_id").Having("max(created_at)").Find(&h)
	return h
}
