package assets

import (
	"embed"
	"encoding/json"
)

//go:embed all:assets
var Assets embed.FS

var DiscordPermissions map[string]any

func init() {
	// Open assets/serenity_perms.json
	f, err := Assets.Open("assets/serenity_perms.json")

	if err != nil {
		panic(err)
	}

	err = json.NewDecoder(f).Decode(&DiscordPermissions)

	if err != nil {
		panic(err)
	}
}
