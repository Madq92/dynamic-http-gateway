package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLastItem(t *testing.T) {
	assert.Equal(t, "c", LastItem("a/b/c", "/"))
	assert.Equal(t, "", LastItem("", "/"))
	assert.Equal(t, "", LastItem("a/b/", "/"))
}

func TestFirstItem(t *testing.T) {
	assert.Equal(t, "a", FirstItem("a/b/c", "/"))
	assert.Equal(t, "", FirstItem("", "/"))
	assert.Equal(t, "", FirstItem("/a/b", "/"))
}

func TestNthItem(t *testing.T) {
	item, err := NthItem("a/b/c/d", "/", 1)
	assert.Equal(t, "b", item)
	assert.Nil(t, err)
	item, err = NthItem("", "/", 1)
	assert.NotNil(t, err)
	item, err = NthItem("/a/b", "/", 3)
	assert.NotNil(t, err)
}
