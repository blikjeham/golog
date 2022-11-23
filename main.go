package main

type buklog struct {
	ui        userInterface
	frequency string
	mode      string
}

var app buklog
var LogChannel = make(chan *LogEntry)

func main() {
	app.frequency = "14.250"
	app.mode = "SSB"
	uiInit()
	uiRun()
}
