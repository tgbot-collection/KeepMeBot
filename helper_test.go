// KeepMeBot - helper_test
// 2020-08-28 20:11
// Benny <benny.think@gmail.com>

package main

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestGetTagAndHash(t *testing.T) {
	var version, hash, message = getTagAndHash()
	assert.NotEmpty(t, version)
	assert.NotEmpty(t, hash)
	assert.NotEmpty(t, message)

}
