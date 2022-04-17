package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/fx/fxtest"
	"go.uber.org/goleak"
)

func TestMain(m *testing.M) {
	goleak.VerifyTestMain(m)
}

func TestOptions(t *testing.T) {
	fxtest.New(t).RequireStart()
	fxtest.New(t).RequireStop()
	assert.True(t, true)
}
