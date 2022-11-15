package main

var frequency string
var mode string

type buklog struct {
	ui userInterface
}

var app buklog

func main() {
	uiInit()
	uiRun()
}
