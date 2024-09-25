package choose_guild

import (
	"context"
	"log/slog"
	"slices"

	"github.com/anti-raid/evil-befall/pkg/api/users"
	"github.com/anti-raid/evil-befall/pkg/state"
	"github.com/anti-raid/evil-befall/pkg/tui"
	"github.com/rivo/tview"
)

type ChooseGuildRoute struct {
	ctx           context.Context
	ctxCancelFunc context.CancelFunc
}

func (r *ChooseGuildRoute) Command() string {
	return "choose_guild"
}

func (r *ChooseGuildRoute) Description() string {
	return "Choose the selected guild"
}

func (r *ChooseGuildRoute) Arguments() [][3]string {
	return [][3]string{
		{"guild_id", "The ID of the guild to choose. If unset, will show guild picker", "string"},
		{"refresh", "Whether to refresh the guild list", "bool"},
	}
}

func (r *ChooseGuildRoute) Setup(state *state.State) error {
	ctx, cancelFunc := context.WithCancel(context.Background())

	r.ctx = ctx
	r.ctxCancelFunc = cancelFunc
	return nil
}

func (r *ChooseGuildRoute) Destroy(state *state.State) error {
	if r.ctxCancelFunc != nil {
		r.ctxCancelFunc()
	}
	return nil
}

func (r *ChooseGuildRoute) Render(state *state.State, args map[string]string) error {
	if guildID, ok := args["guild_id"]; ok {
		return state.SetSelectedGuild(guildID)
	}

	refresh := false

	if refreshStr, ok := args["refresh"]; ok {
		if refreshStr == "true" {
			refresh = true
		}
	}

	slog.Info("Fetching user guild list...", slog.Bool("refresh", refresh))

	guilds, err := users.GetUserGuilds(r.ctx, state, &users.GetUserGuildsData{Refresh: refresh})

	if err != nil {
		return err
	}

	// Create a tview and render it
	type ButtonActionType int

	const (
		ButtonActionSelectGuild   ButtonActionType = iota
		ButtonActionInviteToGuild ButtonActionType = iota
		ButtonActionExit          ButtonActionType = iota
	)

	var continueChan = make(chan ButtonActionType)
	var doneChan = make(chan struct{})

	form := tview.NewForm()

	var dropdownOptions []string

	for _, guild := range guilds.Guilds {
		dropdownOptions = append(dropdownOptions, guild.Name+" ("+guild.ID+")")
	}

	form.AddDropDown("Guilds", dropdownOptions, 0, func(option string, optionIndex int) {
		guildID := guilds.Guilds[optionIndex].ID

		// Check if the bot is in the guild

		var buttonAction ButtonActionType
		if !slices.Contains(guilds.BotInGuilds, guildID) {
			buttonAction = ButtonActionInviteToGuild
		} else {
			buttonAction = ButtonActionSelectGuild
		}

		form.ClearButtons()

		switch buttonAction {
		case ButtonActionSelectGuild:
			form.AddButton("Select Guild", func() {
				continueChan <- ButtonActionSelectGuild
			})
		case ButtonActionInviteToGuild:
			form.AddButton("Invite to Guild", func() {
				continueChan <- ButtonActionInviteToGuild
			})
		}

		form.AddButton("Exit", func() {
			continueChan <- ButtonActionExit
		})
	})

	app := tui.NewTview(state)
	app.SetRoot(form, true)

	go func() {
		for {
			select {
			case v := <-continueChan:
				app.Stop()

				switch v {
				case ButtonActionExit:
					doneChan <- struct{}{}
				case ButtonActionSelectGuild:
					idx, _ := form.GetFormItemByLabel("Guilds").(*tview.DropDown).GetCurrentOption()

					if idx < 0 || idx >= len(guilds.Guilds) {
						slog.Error("Invalid guild index", slog.Int("idx", idx))
						doneChan <- struct{}{}
						return
					}

					if err := state.SetSelectedGuild(guilds.Guilds[idx].ID); err != nil {
						slog.Error("Failed to set selected guild", slog.String("err", err.Error()))
					}

					doneChan <- struct{}{}
				case ButtonActionInviteToGuild:
					idx, _ := form.GetFormItemByLabel("Guilds").(*tview.DropDown).GetCurrentOption()

					if idx < 0 || idx >= len(guilds.Guilds) {
						slog.Error("Invalid guild index", slog.Int("idx", idx))
						doneChan <- struct{}{}
						return
					}

					// TODO: Actually generate an invite link
					slog.Info("You need to invite the bot to the guild", slog.String("guild_id", guilds.Guilds[idx].ID))
					doneChan <- struct{}{}
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

	// Done here
	<-doneChan

	return nil
}
