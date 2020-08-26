// KeepMeBot - database
// 2020-08-20 22:01
// Benny <benny.think@gmail.com>

package main

import (
	"fmt"
	"github.com/jinzhu/gorm"
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

type Session struct {
	BaseModel
	UserID int `gorm:"unique"`
	Next   string
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
	DB.AutoMigrate(&Service{}, &Queue{}, &History{}, &Session{})

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

func addQueue(userid int, username, text, keepType string) (message string) {
	// user id, commands
	s := getServiceMap()
	data := s[keepType]
	format := strings.Count(data.Command, "%s")
	var inputs []interface{}
	for i := 0; i < format; i++ {
		inputs = append(inputs, shellescape.Quote(text)) // not good
	}
	realCommand := fmt.Sprintf(data.Command, inputs...)
	// max than 5?
	count := 0
	DB.Model(&Queue{}).Where("user_id = ? and service_type=?", userid, keepType).Count(&count)

	if count > data.Max {
		message = fmt.Sprintf("Your limit is %d, you are using %d", data.Max, data)
	} else {
		d := Queue{
			UserID:      userid,
			UserName:    username,
			Command:     realCommand,
			ServiceType: keepType,
		}
		DB.Create(&d)
		message = fmt.Sprintf("%s Your command is `%s`", "Success!", realCommand)
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
