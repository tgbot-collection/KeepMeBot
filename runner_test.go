// KeepMeBot - runner_test
// 2020-08-28 20:11
// Benny <benny.think@gmail.com>

package main

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestGet(t *testing.T) {
	actual := get("http://z.cn")
	expected := fmt.Sprintf("http status code is %d", 200)
	assert.Equal(t, expected, actual)
}

func TestScheduler(t *testing.T) {
	// two queue, three history
	q1 := Queue{
		UserID:    8888,
		ServiceID: 1,
	}
	q2 := Queue{
		UserID:    8888,
		Command:   "http://baidu.com",
		ServiceID: 3,
	}
	DB.Unscoped().Delete(&Queue{})
	DB.Unscoped().Delete(&History{})
	DB.Create(&q1)
	DB.Create(&q2)

	scheduler()

	h := getNewestHistory(8888)
	assert.Equal(t, 2, len(h))
	assert.Equal(t, h[0].ServiceID, 1)
	assert.Equal(t, h[1].ServiceID, 3)
	time.Sleep(time.Second * 12)
	scheduler()
	DB.Model(&History{}).Where("user_id=?", 8888).Find(&h)
	assert.Equal(t, 3, len(h))

}

func TestExternalExecutor(t *testing.T) {
	q1 := Queue{
		Command: "echo 123",
	}
	q2 := Queue{
		Command: "gfdda",
	}
	resp := externalExecutor(q1)
	assert.Contains(t, resp, "123")
	resp = externalExecutor(q2)
	assert.Contains(t, resp, "not found")

}

func TestInternalExecutor(t *testing.T) {
	q1 := Queue{
		Command:   "get http://www.qq.com/",
		Parameter: "http://www.qq.com/",
		ServiceID: 3,
	}
	resp := internalExecutor(q1)
	assert.Contains(t, resp, "http status code")
}
