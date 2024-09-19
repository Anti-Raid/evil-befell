package login

import (
	"context"
	"strconv"

	"github.com/anti-raid/evil-befall/pkg/api/core"
	"github.com/anti-raid/evil-befall/pkg/auth"
	"github.com/anti-raid/evil-befall/pkg/constants"
	"github.com/anti-raid/evil-befall/pkg/state"
	"github.com/anti-raid/evil-befall/pkg/statusbox"
	"github.com/pkg/browser"
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

		state.StateFetchOptions.InstanceAPIUrl = instanceUrl

		loginStatusWriter.AddStatusMessage("Logging in... | Attempt #" + strconv.Itoa(i))
		i++

		apiConfig, err := core.GetApiConfig(r.ctx, state)

		if err != nil {
			loginStatusWriter.AddStatusMessage("Failed to get API config: " + err.Error())
			return
		}

		loginUrl := auth.GetAuthURL(r.ctx, state, apiConfig)

		loginStatusWriter.AddStatusMessage("Login URL: " + loginUrl)

		if err := browser.OpenURL(loginUrl); err != nil {
			loginStatusWriter.AddStatusMessage("Failed to open browser: " + err.Error())
		}

		ul, err := auth.CreateSessionOnServerAddr(r.ctx, state)

		if err != nil {
			loginStatusWriter.AddStatusMessage("Failed to create session: " + err.Error())
			return
		}

		loginStatusWriter.AddStatusMessage("Session created: " + ul.UserID)
	})

	form.AddButton("Exit", func() {
		app.Stop()
	})

	// Add the form to the grid
	grid.AddItem(form, 0, 0, 1, 3, 0, 0, true)

	// Add the status box to the grid below the form taking up the remaining screen width and height
	grid.AddItem(loginStatusBox, 1, 0, 1, 3, 0, 0, false)

	return grid, nil
}
