package tui

import (
	"github.com/anti-raid/evil-befall/pkg/state"
	"github.com/rivo/tview"
)

func NewTview(state *state.State) *tview.Application {
	return tview.NewApplication().EnableMouse(state.Prefs.MouseEnabledInTView).EnablePaste(state.Prefs.PasteEnabledInTView)
}
