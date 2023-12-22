package main

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewAccount(t *testing.T) {
	account, err := NewAccount("John", "Doe", "password")

	assert.Nil(t, err)

	fmt.Printf("Account: %+v\n", account)
}
