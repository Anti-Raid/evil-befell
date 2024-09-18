package main

import (
	"os"

	"github.com/anti-raid/evil-befall/pkg/router"
	_ "github.com/anti-raid/evil-befall/pkg/routes"
	"github.com/anti-raid/evil-befall/pkg/state"
	"github.com/rivo/tview"
)

func envOrBool(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}

	return fallback
}

func main() {
	// Create a new state
	var state = state.NewState()

	var mouseEnabled = envOrBool("MOUSE_ENABLED", "true") == "true"
	var pasteEnabled = envOrBool("PASTE_ENABLED", "true") == "true"
	var fullscreen = envOrBool("FULLSCREEN", "true") == "true"

	var app = tview.NewApplication().EnableMouse(mouseEnabled).EnablePaste(pasteEnabled)
	var pages = tview.NewPages()

	_, err := router.GotoCurrent(state, app, pages)

	if err != nil {
		panic(err)
	}

	if err := app.SetRoot(pages, fullscreen).Run(); err != nil {
		panic(err)
	}
}
