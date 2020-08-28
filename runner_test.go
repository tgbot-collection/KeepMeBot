// KeepMeBot - runner_test
// 2020-08-28 20:11
// Benny <benny.think@gmail.com>

package main

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestGet(t *testing.T) {
	actual := get("http://z.cn")
	expected := fmt.Sprintf("http status code is %d", 200)
	assert.Equal(t, expected, actual)
}

func TestInternalExecutor(t *testing.T) {

}
