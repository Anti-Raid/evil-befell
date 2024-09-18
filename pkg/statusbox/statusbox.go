package statusbox

import (
	"context"

	"github.com/rivo/tview"
)

// A status box is a primitive that displays a status message.
type StatusBox struct {
	ctx      context.Context
	messages chan string
	w        *tview.TextView
}

// NewStatusBox creates a new status box.
func NewStatusBox(ctx context.Context, tv *tview.TextView) *StatusBox {
	sb := &StatusBox{
		ctx:      ctx,
		messages: make(chan string, 100),
		w:        tv,
	}

	go sb.drawMessages()

	return sb
}

// Add a new status message to the status box.
func (s *StatusBox) AddStatusMessage(msg string) {
	s.messages <- msg
}

func (r *StatusBox) drawMessages() {
	defer recover()

	// Draw the message
	for {
		select {
		// Context closed case
		case <-r.ctx.Done():
			return
		default:
			continue
		case msg := <-r.messages:
			r.w.Write([]byte(msg + "\n"))

			r.w.ScrollToEnd()
		}
	}
}
