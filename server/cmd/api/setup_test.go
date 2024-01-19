package main

import (
	"os"
	"testing"
)

var app Application

func TestMain(m *testing.M) {
	os.Exit(m.Run())
}
