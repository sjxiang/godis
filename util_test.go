package main

import (
	"testing"
)

// telnet 127.0.0.1 5001

func TestTelnet(t *testing.T) {
	LogMessage(Red,     "1")
	LogMessage(Blue,    "2")
	LogMessage(Green,   "3")
	LogMessage(Yellow,  "4")
	LogMessage(Magenta, "5")
	LogMessage(Cyan,    "6")
	LogMessage(White,   "7")
}