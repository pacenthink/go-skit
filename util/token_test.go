package util

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_ParseBearerJwtFromAuthHeader(t *testing.T) {
	_, err := ParseBearerJwtFromAuthHeader("fobbar")
	assert.NotNil(t, err)

	_, err = ParseBearerJwtFromAuthHeader("blubber fobbar")
	assert.NotNil(t, err)

	_, err = ParseBearerJwtFromAuthHeader("bearer fobbar")
	assert.NotNil(t, err)

	_, err = ParseBearerJwtFromAuthHeader("bearer ")
	assert.NotNil(t, err)
}
