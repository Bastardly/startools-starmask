package ui

import (
	"strconv"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/driver/mobile"
	"fyne.io/fyne/v2/widget"
)

type numericalEntry struct {
	widget.Entry
}

func newNumericalEntry() *numericalEntry {
	entry := &numericalEntry{}
	entry.ExtendBaseWidget(entry)
	return entry
}

func (e *numericalEntry) setDot(r rune) {
	hasdot := false
	for i := 0; i < len(e.Entry.Text); i++ {
		if e.Entry.Text[i] == '.' || e.Entry.Text[i] == ',' {
			hasdot = true
			break
		}
	}
	if !hasdot {
		e.Entry.TypedRune('.')
	}
}

func (e *numericalEntry) TypedRune(r rune) {
	switch r {
	case '.', ',':
		e.setDot(r)

	case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
		e.Entry.TypedRune(r)
	}
}

func (e *numericalEntry) TypedShortcut(shortcut fyne.Shortcut) {
	paste, ok := shortcut.(*fyne.ShortcutPaste)
	if !ok {
		e.Entry.TypedShortcut(shortcut)
		return
	}

	content := paste.Clipboard.Content()
	if _, err := strconv.ParseFloat(content, 64); err == nil {
		e.Entry.TypedShortcut(shortcut)
	}
}

func (e *numericalEntry) Keyboard() mobile.KeyboardType {
	return mobile.NumberKeyboard
}

func numericEntry() *numericalEntry {

	return newNumericalEntry()
}
