package login

import (
	"context"
	"strconv"

	"github.com/anti-raid/evil-befall/pkg/constants"
	"github.com/anti-raid/evil-befall/pkg/state"
	"github.com/anti-raid/evil-befall/pkg/statusbox"
	"github.com/rivo/tview"
)

type LoginRoute struct {
	ctx           context.Context
	ctxCancelFunc context.CancelFunc
}

func (r *LoginRoute) ID() string {
	return "login"
}

func (r *LoginRoute) Title() string {
	return "Root"
}

func (r *LoginRoute) Description() string {
	return "The root route"
}

func (r *LoginRoute) Setup(state *state.State) error {
	ctx, cancelFunc := context.WithCancel(context.Background())

	r.ctx = ctx
	r.ctxCancelFunc = cancelFunc
	return nil
}

func (r *LoginRoute) Destroy(state *state.State) error {
	r.ctxCancelFunc()
	return nil
}

func (r *LoginRoute) Render(state *state.State, app *tview.Application, pages *tview.Pages) (tview.Primitive, error) {
	loginStatusBox := tview.NewTextView()
	loginStatusWriter := statusbox.NewStatusBox(r.ctx, loginStatusBox)

	grid := tview.NewGrid().SetColumns(0, 0, 0).SetRows(0, 0, 0)

	form := tview.NewForm()

	// Create a text box prompting for instance URL
	form.AddInputField("Instance URL", constants.DefaultInstanceUrl, 0, nil, nil)

	// Create a new button for login
	i := 0
	form.AddButton("Login", func() {
		instanceUrl := form.GetFormItem(0).(*tview.InputField).GetText()
		if instanceUrl == "" {
			panic("Instance URL is required")
		}

		loginStatusWriter.AddStatusMessage("Logging in..." + strconv.Itoa(i))
		i++
	})

	form.AddButton("Exit", func() {
		app.Stop()
	})

	// Add the form to the grid
	grid.AddItem(form, 0, 0, 1, 3, 0, 0, true)

	// Add the status box to the grid
	grid.AddItem(loginStatusBox, 1, 0, 1, 3, 0, 0, false)

	return grid, nil
}
