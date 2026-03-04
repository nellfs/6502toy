package terminal

import (
	"github.com/diamondburned/gotk4/pkg/core/glib"
	"github.com/diamondburned/gotk4/pkg/gtk/v4"
)

type TextViewWriter struct {
	view    *gtk.TextView
	buffer  *gtk.TextBuffer
	endMark *gtk.TextMark
}

func NewTextViewWriter(view *gtk.TextView) *TextViewWriter {
	buffer := view.Buffer()
	// create one persistent mark at the end, gravity=false means it stays at end
	end := buffer.EndIter()
	endMark := buffer.CreateMark("terminal-end", end, false)

	return &TextViewWriter{
		view:    view,
		buffer:  buffer,
		endMark: endMark,
	}
}

func (w *TextViewWriter) Write(p []byte) (int, error) {
	text := string(p)
	glib.IdleAdd(func() {
		end := w.buffer.EndIter()
		w.buffer.Insert(end, text)
		w.buffer.MoveMarkByName("terminal-end", w.buffer.EndIter())

		// second idle: scroll AFTER GTK has remeasured the new content height
		glib.IdleAdd(func() {
			w.view.ScrollToMark(w.endMark, 0.0, true, 0.0, 1.0)
		})
	})
	return len(p), nil
}
