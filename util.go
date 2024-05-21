package main

import (
	"fmt"
)

type Color string

const (
	Red     Color = "red"
	Blue    Color = "blue"
	Green   Color = "green"
	Yellow  Color = "yellow"
	Magenta Color = "magenta"  
	Cyan    Color = "cyan"     
	White   Color = "white"
)

func LogMessage(color Color, message string) {

	switch color {
	case Red:
		fmt.Printf("\033[0;31m%s\033[0m\n", message)
	case Green:
		fmt.Printf("\033[0;32m%s\033[0m\n", message)
	case Yellow:
		fmt.Printf("\033[0;33m%s\033[0m\n", message)
	case Blue:
		fmt.Printf("\033[0;34m%s\033[0m\n", message)
	case Magenta:
		fmt.Printf("\033[0;35m%s\033[0m\n", message)
	case Cyan:
		fmt.Printf("\033[0;36m%s\033[0m\n", message)
	case White:
		fmt.Printf("\033[0;37m%s\033[0m\n", message)
	default:
		fmt.Println(message)
	}
}
