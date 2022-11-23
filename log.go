package main

import (
	"time"
)

type LogEntry struct {
	datetime  time.Time
	callsign  string
	frequency string
	mode      string
	rst_rx    string
	rst_tx    string
}

func NewLogEntry(callsign, frequency, mode, rst_rx, rst_tx string) *LogEntry {
	entry := &LogEntry{
		callsign:  callsign,
		frequency: frequency,
		mode:      mode,
		rst_rx:    rst_rx,
		rst_tx:    rst_tx,
	}
	entry.datetime = time.Now().UTC()
	return entry
}

func (entry *LogEntry) Store() {
	// Store entry in database
	LogChannel <- entry
}
