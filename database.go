// KeepMeBot - database
// 2020-08-20 22:01
// Benny <benny.think@gmail.com>

package main

import (
	"fmt"
	"github.com/jinzhu/gorm"
	"gopkg.in/tucnak/telebot.v2"
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
	Command     string
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

func init() {
	supportedService = []Service{
		{
			Name:        "Docker Hub",
			Max:         5,
			ServiceType: "external",
			Command:     "docker pull %s && docker rmi %s",
		},
		{
			Name:        "GitHub",
			Max:         3,
			ServiceType: "external",
			Command:     "git clone %s && rm -rf %s",
		},
		{
			Name:        "get",
			Max:         10,
			ServiceType: "internal",
			Command:     "get %s",
		},
	}
	var err error
	DB, err = gorm.Open("sqlite3", "keep.db")
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

	realCommand := fmt.Sprintf(service.Command, inputs...)
	// max than 5
	count := 0
	query := Queue{
		UserID: m.Sender.ID}
	DB.Model(&query).Count(&count)
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
	}
	DB.Create(&h)
}
