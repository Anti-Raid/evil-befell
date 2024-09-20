package main

import (
	"flag"
	"fmt"
	"log/slog"
	"os"

	"github.com/anti-raid/evil-befall/pkg/router"
	_ "github.com/anti-raid/evil-befall/pkg/routes"
	statelib "github.com/anti-raid/evil-befall/pkg/state"
	"github.com/infinitybotlist/eureka/shellcli"
)

type cliData struct {
	State *statelib.State
}

func envOrBool(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}

	return fallback
}

func envOrString(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}

	return fallback
}

func main() {
	// Create a new state
	var mouseEnabled = envOrBool("MOUSE_ENABLED", "false") == "true"
	var pasteEnabled = envOrBool("PASTE_ENABLED", "true") == "true"
	var fullscreen = envOrBool("FULLSCREEN", "true") == "true"
	var persist = envOrString("PERSIST", "evil-befall-cfg.json")

	// Set state.Prefs
	state, err := statelib.NewState(statelib.UserPref{
		MouseEnabledInTView:      mouseEnabled,
		PasteEnabledInTView:      pasteEnabled,
		FullscreenEnabledInTView: fullscreen,
		Persist: func() *string {
			if persist == "" || persist == "false" {
				return nil
			}

			return &persist
		}(),
	})

	if err != nil {
		slog.Error("Failed to create state:", slog.String("error", err.Error()))
		os.Exit(1)
	}

	// Create command list
	var commands = make(map[string]*shellcli.Command[cliData])

	for _, route := range router.Routes() {
		commands[route.Command()] = &shellcli.Command[cliData]{
			Description: route.Description(),
			Args:        route.Arguments(),
			Run: func(cli *shellcli.ShellCli[cliData], args map[string]string) error {
				return router.Goto(route.Command(), cli.Data.State, args)
			},
		}
	}

	root := &shellcli.ShellCli[cliData]{
		Data: &cliData{
			State: state,
		},
		Prompter: func(r *shellcli.ShellCli[cliData]) string {
			return "evil-befall> "
		},
		Commands: commands,
	}

	root.AddCommand("help", root.Help())

	// Handle --command args
	command := flag.String("command", "", "Command to run. If unset, will run as shell")
	flag.Parse()

	if command != nil && *command != "" {
		err := root.Init()

		if err != nil {
			fmt.Println("Error initializing cli: ", err)
			os.Exit(1)
		}

		cancel, err := root.ExecuteCommands(*command)

		if err != nil {
			fmt.Println("Error:", err)
		}

		if cancel {
			fmt.Println("Exiting...")
		}

		return
	}

	root.Run()
}
