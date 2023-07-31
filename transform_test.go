package main

import (
	"github.com/stretchr/testify/assert"
	"strings"
	"testing"
)

func TestHelloWorld(t *testing.T) {
	requests := ReadMapping()
	assert.GreaterOrEqual(t, len(requests), 1)

	response := ExpressRule(requests, &Request{
		Host:   "myhost.com",
		Path:   "/aa.php",
		Method: "GET",
	})

	assert.Equal(t, true, strings.Contains(response, "properties"))
}
