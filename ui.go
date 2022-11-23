package main

import (
	"fmt"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	"strings"
	"time"
)

type userInterface struct {
	app    *tview.Application
	pages  *tview.Pages
	header *tview.Grid
	form   *tview.Form
	log    *tview.Table
	focus  string
}

var timeFmt = "2006-01-02T15:04:05"

func uiPeriodic() {
	for {
		time.Sleep(time.Second / 2)
		updateHeader(app.ui.header)
		app.ui.app.Draw()
	}
}

func uiInit() {
	app.ui.app = tview.NewApplication()
	makePages()
	go uiPeriodic()
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
	main_window.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Key() == tcell.KeyF1 {
			app.ui.pages.ShowPage("help")
			return nil
		}
		if event.Key() == tcell.KeyF2 {
			app.ui.pages.ShowPage("freq")
			return nil
		}
		if event.Key() == tcell.KeyF3 {
			switch app.ui.focus {
			case "form":
				app.ui.app.SetFocus(app.ui.log)
				app.ui.focus = "log"
			case "log":
				app.ui.app.SetFocus(app.ui.form)
				app.ui.focus = "form"
			}
			return nil
		}
		return event
	})
}

func makeMainWindow() *tview.Grid {
	grid := makeGrid()
	app.ui.form = makeForm()
	app.ui.log = makeLog()
	grid.AddItem(app.ui.form, 1, 0, 1, 1, 0, 0, true)
	grid.AddItem(app.ui.log, 1, 1, 1, 2, 0, 0, false)
	app.ui.focus = "form"
	return grid
}

func makeGrid() *tview.Grid {

	grid := tview.NewGrid().
		SetRows(3, 0, 1).
		SetColumns(30, 0).
		SetBorders(true)

	header := makeHeader()
	footer := "F1 help | F2 frequency | F3 list | Ctrl-C exit"
	grid.AddItem(header, 0, 0, 1, 3, 0, 0, false)
	grid.AddItem(make_text(footer), 2, 0, 1, 3, 0, 0, false)

	return grid
}

func updateFooter(footer *tview.Grid) {

}

func updateHeader(header *tview.Grid) {
	header.Clear()
	header.AddItem(center_text("PA5BUK BUKLog"), 0, 0, 1, 3, 0, 0, false)
	freq_header := fmt.Sprintf("frequency: %s; mode: %s", app.frequency, app.mode)
	header.AddItem(make_text(freq_header), 1, 0, 1, 2, 0, 0, false)

	time_s := fmt.Sprintf(time.Now().UTC().Format(timeFmt))
	header.AddItem(center_text(time_s), 1, 2, 1, 1, 0, 0, false)
}

func makeHeader() *tview.Grid {
	header := tview.NewGrid().
		SetRows(0, 1).
		SetColumns(-3, 5).
		SetBorders(false)

	updateHeader(header)
	app.ui.header = header
	return header
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
	form.AddButton("Save", func() {
		c := strings.ToUpper(callsign.GetText())
		t := rst_tx.GetText()
		r := rst_rx.GetText()
		entry := NewLogEntry(
			c,
			app.frequency,
			app.mode,
			t,
			r,
		)
		entry.Store()
		// addLogEntry(app.ui.log, entry)

		// Clear form
		callsign.SetText("")
		rst_tx.SetText("")
		rst_rx.SetText("")
		form.SetFocus(1)
	})
	return form
}

func make_header(table *tview.Table) {
	table.SetCellSimple(0, 0, "Date/time")
	table.SetCellSimple(0, 1, "Callsign")
	table.SetCellSimple(0, 2, "Frequency")
	table.SetCellSimple(0, 3, "Mode")
	table.SetFixed(0, 1)
}

func makeRow(table *tview.Table, r int, entry *LogEntry) {
	table.SetCell(r, 0, tview.NewTableCell(entry.datetime.Format(timeFmt)))
	table.SetCell(r, 1, tview.NewTableCell(entry.callsign))
	table.SetCell(r, 2, tview.NewTableCell(entry.frequency))
	table.SetCell(r, 3, tview.NewTableCell(entry.mode))
}

func addLogEntry(table *tview.Table, entry *LogEntry) {
	row := table.GetRowCount()
	makeRow(table, row, entry)
}

func entryReceiver() {
	for {
		entry := <-LogChannel
		addLogEntry(app.ui.log, entry)
	}
}

func makeLog() *tview.Table {
	var le LogEntry
	table := tview.NewTable().SetBorders(true).SetFixed(1, 0)
	app.ui.log = table
	make_header(table)
	go entryReceiver()

	// Add dummy entry
	le.datetime = time.Now().UTC()
	le.callsign = "PA5BUK"
	le.frequency = "14.3145"
	le.mode = "SSB"
	LogChannel <- &le

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
	f := tview.NewInputField().SetLabel("Frequency")
	m := tview.NewDropDown().SetLabel("Mode").SetOptions([]string{"SSB", "CW", "AM", "FM", "Digi"}, nil)
	form.AddFormItem(f)
	form.AddFormItem(m)
	form.AddButton("Save", func() {
		app.frequency = f.GetText()
		_, app.mode = m.GetCurrentOption()
		updateHeader(app.ui.header)
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
