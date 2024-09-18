package root

import (
	"github.com/anti-raid/evil-befall/pkg/router"
	"github.com/anti-raid/evil-befall/pkg/state"
	"github.com/rivo/tview"
)

type RootRoute struct{}

func (r *RootRoute) ID() string {
	return "root"
}

func (r *RootRoute) Title() string {
	return "Root"
}

func (r *RootRoute) Description() string {
	return "The root route"
}

func (r *RootRoute) Setup(state *state.State) error {
	return nil
}

func (r *RootRoute) Destroy(state *state.State) error {
	return nil
}

func (r *RootRoute) Render(state *state.State, app *tview.Application, pages *tview.Pages) (tview.Primitive, error) {
	form := tview.NewForm()

	// Create a new button for login
	form.AddButton("Login", func() {
		_, err := router.Goto(state, "login", app, pages)

		if err != nil {
			panic(err)
		}
	})
	form.AddTextView("Click the login button to login", "", 0, 1, true, true)
	form.AddButton("Exit", func() {
		app.Stop()
	})

	return form, nil
}
