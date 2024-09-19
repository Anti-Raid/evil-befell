package login

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/anti-raid/evil-befall/pkg/api/core"
	"github.com/anti-raid/evil-befall/pkg/auth"
	"github.com/anti-raid/evil-befall/pkg/constants"
	"github.com/anti-raid/evil-befall/pkg/state"
	"github.com/anti-raid/evil-befall/pkg/tui"
	"github.com/pkg/browser"
	"github.com/rivo/tview"
)

type LoginRoute struct {
	ctx           context.Context
	ctxCancelFunc context.CancelFunc
}

func (r *LoginRoute) Command() string {
	return "login"
}

func (r *LoginRoute) Description() string {
	return "Log in to an Anti-Raid instance"
}

func (r *LoginRoute) Arguments() [][3]string {
	return [][3]string{}
}

func (r *LoginRoute) Setup(state *state.State) error {
	ctx, cancelFunc := context.WithCancel(context.Background())

	r.ctx = ctx
	r.ctxCancelFunc = cancelFunc
	return nil
}

func (r *LoginRoute) Destroy(state *state.State) error {
	if r.ctxCancelFunc != nil {
		r.ctxCancelFunc()
	}
	return nil
}

func (r *LoginRoute) Render(state *state.State, args map[string]string) error {
	var continueChan = make(chan bool)
	var doneChan = make(chan struct{})

	form := tview.NewForm()

	// Create a text box prompting for instance URL
	form.AddInputField("Instance URL", constants.DefaultInstanceUrl, 0, nil, nil)

	// Create a new button for login
	form.AddButton("Login", func() {
		continueChan <- true
	})

	form.AddButton("Exit", func() {
		continueChan <- false
	})

	app := tui.NewTview(state)
	app.SetRoot(form, true)

	go func() {
		for {
			select {
			case v := <-continueChan:
				app.Stop()
				if v {
					instanceUrl := form.GetFormItemByLabel("Instance URL").(*tview.InputField).GetText()
					state.StateFetchOptions.InstanceAPIUrl = instanceUrl
					err := execLogin(r, state)
					if err != nil {
						slog.Error("Failed to login", slog.String("err", err.Error()))
					}
				}

				doneChan <- struct{}{}
			case <-r.ctx.Done():
				app.Stop()
				doneChan <- struct{}{}
			}
		}
	}()

	if err := app.Run(); err != nil {
		return err
	}

	<-doneChan

	return nil
}

func execLogin(r *LoginRoute, state *state.State) error {
	slog.Info("Fetching API config", slog.String("instanceUrl", state.StateFetchOptions.InstanceAPIUrl))

	apiConfig, err := core.GetApiConfig(r.ctx, state)

	if err != nil {
		return err
	}

	loginUrl := auth.GetAuthURL(r.ctx, state, apiConfig)

	slog.Info("Please visit the following url and login", slog.String("loginUrl", loginUrl))

	if err := browser.OpenURL(loginUrl); err != nil {
		slog.Error("Failed to open browser", slog.String("err", err.Error()))
	}

	ul, err := auth.CreateSessionOnServerAddr(r.ctx, state)

	if err != nil {
		return fmt.Errorf("failed to create session: %w", err)
	}

	slog.Info("Session created", slog.String("userId", ul.UserID), slog.String("sessionId", ul.SessionID))

	state.Session.AddSession(ul)

	return nil
}
