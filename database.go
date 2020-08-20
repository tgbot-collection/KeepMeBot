// KeepMeBot - database
// 2020-08-20 22:01
// Benny <benny.think@gmail.com>

package main

import (
	"fmt"
	"github.com/jinzhu/gorm"
	tb "gopkg.in/tucnak/telebot.v2"
	"strings"
)

import "gopkg.in/alessio/shellescape.v1"

// SQLite

// all the support keep alive service are bundle with database.

var DB *gorm.DB

type BaseModel struct {
	gorm.Model
}

type Service struct {
	BaseModel
	Name    string `gorm:"unique"`
	Max     int
	Command string
}

type Queue struct {
	BaseModel
	UserID      int
	UserName    string
	Command     string
	ServiceType string
}

type History struct {
	BaseModel
	UserID   int
	UserName string
	Command  string
	Output   string
}

func init() {
	var supportedService = []Service{
		{
			Name:    "Docker Hub",
			Max:     5,
			Command: "docker pull %s && docker rmi %s",
		},
		{
			Name:    "GitHub",
			Max:     3,
			Command: "git clone %s && rm -rf %s",
		},
	}
	var err error
	DB, err = gorm.Open("sqlite3", "keep.DB")
	if err != nil {
		panic("failed to connect database")
	}

	// Migrate the schema
	DB.AutoMigrate(&Service{}, &Queue{}, &History{})

	// 创建
	DB.Unscoped().Delete(&Service{})
	for _, v := range supportedService {
		DB.Create(&v)
	}

}

func getServiceArray() (supportedService []Service, err error) {
	if err = DB.Find(&supportedService).Error; err != nil {
		return
	}
	return
}
func getServiceMap() map[string]Service {
	var a = make(map[string]Service)
	arr, _ := getServiceArray()

	for _, v := range arr {
		a[v.Name] = v
	}
	return a
}

func addQueue(m *tb.Message, keepType string) (status string, realCommand string) {
	// user id, commands
	s := getServiceMap()
	data := s[keepType]
	format := strings.Count(data.Command, "%s")
	var inputs []interface{}
	for i := 0; i < format; i++ {
		inputs = append(inputs, shellescape.Quote(m.Text)) // not good
	}
	realCommand = fmt.Sprintf(data.Command, inputs...)
	// max than 5?
	count := 0
	DB.Model(&Queue{}).Where("user_id = ? and service_type=?", m.Sender.ID, keepType).Count(&count)

	if count > data.Max {
		status = fmt.Sprintf("Your limit is %d, you are using %d", data.Max, data)
	} else {
		d := Queue{
			BaseModel:   BaseModel{},
			UserID:      m.Sender.ID,
			UserName:    m.Sender.Username,
			Command:     realCommand,
			ServiceType: keepType,
		}
		DB.Create(&d)
		status = "Success!"
	}
	DB.Find(&Queue{})
	return
}
