package main

import (
	"fmt"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

type userInterface struct {
	app *tview.Application
	pages *tview.Pages
}

func uiInit() {
	app.ui.app = tview.NewApplication()
	makePages()
}

func uiRun() {
	if err := app.ui.app.Run(); err != nil {
		panic(err)
	}
}

func makePages() {
	main_window := makeMainWindow()
	help_window := makeHelpWindow()
	freq_window := makeFreqWindow()

	app.ui.pages = tview.NewPages().
		AddPage("main", main_window, true, true).
		AddPage("help", help_window, true, false).
		AddPage("freq", freq_window, true, false)
		help_window.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Key() == tcell.KeyEscape {
			app.ui.pages.SwitchToPage("main")
			return nil
		}
		return event
	})
	freq_window.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Key() == tcell.KeyEscape {
			app.ui.pages.SwitchToPage("main")
			return nil
		}
		return event
	})
	
	app.ui.app.SetRoot(app.ui.pages, true)
	app.ui.app.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Key() == tcell.KeyF1 {
			app.ui.pages.ShowPage("help")
			return nil
		}
		if event.Key() == tcell.KeyF2 {
			app.ui.pages.ShowPage("freq")
			return nil
		}
		return event
	})
}

func makeMainWindow() *tview.Grid {
	grid := makeGrid()
	form := makeForm()
	log := makeLog()
	grid.AddItem(form, 1, 0, 1, 1, 0, 0, true)
	grid.AddItem(log, 1, 1, 1, 2, 0, 0, false)
	return grid
}

func makeGrid() *tview.Grid {

	grid := tview.NewGrid().
		SetRows(3, 0, 1).
		SetColumns(30, 0).
		SetBorders(true)

	header := tview.NewGrid().
		SetRows(0, 1).
		SetColumns(-3, 5).
		SetBorders(false)

	freq_header := fmt.Sprintf("frequency: %s; mode: %s", frequency, mode)
	footer := "F1 help | F2 frequency | F3 list | Ctrl-C exit"

	header.AddItem(center_text("PA5BUK BUKLog"), 0, 0, 1, 3, 0, 0, false)
	header.AddItem(make_text(freq_header), 1, 0, 1, 2, 0, 0, false)
	header.AddItem(center_text("2022-11-14T20:44:15"), 1, 2, 1, 1, 0, 0, false)
	grid.AddItem(header, 0, 0, 1, 3, 0, 0, false)
	grid.AddItem(make_text(footer), 2, 0, 1, 3, 0, 0, false)

	return grid
}

func makeForm() *tview.Form {
	form := tview.NewForm()
	form.SetBorder(true)
	form.SetTitle("Enter log")
	callsign := tview.NewInputField().
		SetLabel("Callsign").
		SetFieldWidth(10)

	rst_tx := tview.NewInputField().
		SetLabel("RST TX").
		SetFieldWidth(6)
	rst_rx := tview.NewInputField().
		SetLabel("RST RX").
		SetFieldWidth(6)
	form.AddFormItem(callsign)
	form.AddFormItem(rst_tx)
	form.AddFormItem(rst_rx)
	form.AddButton("Save", nil)
	return form
}

func make_header(table *tview.Table) {
	table.SetCellSimple(0, 0, "Date/time")
	table.SetCellSimple(0, 1, "Callsign")
	table.SetCellSimple(0, 2, "Frequency")
	table.SetCellSimple(0, 3, "Mode")
	table.SetFixed(0, 1)
}

func make_row(table *tview.Table, r int, datetime, callsign string) {
	table.SetCell(r, 0, tview.NewTableCell(datetime))
	table.SetCell(r, 1, tview.NewTableCell(callsign))
	table.SetCell(r, 2, tview.NewTableCell("14.3145"))
	table.SetCell(r, 3, tview.NewTableCell("SSB"))
}

func makeLog() *tview.Table {
	var le LogEntry
	table := tview.NewTable().SetBorders(true)
	le.datetime = "2022-11-03T13:37"
	le.callsign = "PA5BUK"
	le.frequency = 14.3145
	le.mode = "SSB"
	make_header(table)
	make_row(table, 1, le.datetime, le.callsign)
	make_row(table, 2, "2022-11-03T13:40", "PD3BRT")
	make_row(table, 3, "2022-11-03T13:58", "PA3GFJ")
	return table
}

func makeHelpWindow() *tview.TextView {
	view := tview.NewTextView().
		SetText("This is the help window")
	return view
}

func makeFreqWindow() *tview.Form {
	form := tview.NewForm()
	form.SetBorder(true)
	form.SetTitle("Change frequency")
	freq := tview.NewInputField().SetLabel("Frequency")
	form.AddFormItem(freq)
	form.AddDropDown("Mode", []string{"SSB", "CW", "AM", "FM", "Digi"}, 0, nil)
	form.AddButton("Save", func() {
		frequency = freq.GetText()
		app.ui.pages.SwitchToPage("main")
	})
	return form
}

func make_text(text string) tview.Primitive {
	return tview.NewTextView().
		SetText(text)
}

func center_text(text string) tview.Primitive {
	return tview.NewTextView().
		SetTextAlign(tview.AlignCenter).
		SetText(text)
}

