package terminal

import (
	"github.com/diamondburned/gotk4/pkg/core/glib"
	"github.com/diamondburned/gotk4/pkg/gtk/v4"
)

type TextViewWriter struct {
	view   *gtk.TextView
	buffer *gtk.TextBuffer
}

func NewTextViewWriter(view *gtk.TextView) *TextViewWriter {
	return &TextViewWriter{
		view:   view,
		buffer: view.Buffer(),
	}
}

func (w *TextViewWriter) Write(p []byte) (int, error) {
	text := string(p)

	glib.IdleAdd(func() {
		end := w.buffer.EndIter()
		w.buffer.Insert(end, text)

		mark := w.buffer.CreateMark("", w.buffer.EndIter(), false)
		w.view.ScrollToMark(mark, 0.0, true, 0.0, 1.0)
	})

	return len(p), nil
}
