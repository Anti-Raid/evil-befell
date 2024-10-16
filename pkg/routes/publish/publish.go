package publish

import (
	"context"
	"errors"
	"fmt"

	"github.com/anti-raid/evil-befall/pkg/api/guilds"
	"github.com/anti-raid/evil-befall/pkg/state"
	"github.com/anti-raid/evil-befall/types"
	orderedmap "github.com/wk8/go-ordered-map/v2"
)

type PublishRoute struct {
}

func (r *PublishRoute) Command() string {
	return "publish.template"
}

func (r *PublishRoute) Description() string {
	return "Publish a template to a setting. Note that other settings may work but this is only intended for use in templates"
}

func (r *PublishRoute) Arguments() [][3]string {
	return [][3]string{
		{"guildId", "The guild id", "[selected guild id]"},
		{"module", "The module", ""},
		{"setting", "The setting ID"},
		{"pkey", "The primary key to use for update", ""},
		{"pkeyValue", "The primary key value to use for update", ""},
		{"key", "The key to update"},
		{"value", "The value to set"},
	}
}

func (r *PublishRoute) Setup(state *state.State) error {
	return nil
}

func (r *PublishRoute) Destroy(state *state.State) error {
	return nil
}

func (r *PublishRoute) Render(state *state.State, args map[string]string) error {
	guildId := state.SelectedOptions.GuildID

	if v, ok := args["guildId"]; ok {
		guildId = v
	}

	module, ok := args["module"]

	if !ok {
		return errors.New("module is required")
	}

	setting, ok := args["setting"]

	if !ok {
		return errors.New("setting is required")
	}

	pkey, ok := args["pkey"]

	if !ok {
		return errors.New("pkey is required")
	}

	pvalue, ok := args["pkeyValue"]

	if !ok {
		return errors.New("pkeyValue is required")
	}

	key, ok := args["key"]

	if !ok {
		return errors.New("key is required")
	}

	value, ok := args["value"]

	if !ok {
		return errors.New("value is required")
	}

	fields := orderedmap.New[string, any]()

	fields.Set(key, value)
	fields.Set(pkey, pvalue)

	resp, err := guilds.SettingsExecute(context.Background(), state, &guilds.SettingsExecuteData{
		GuildID: guildId,
		SettingsExecuteData: &types.SettingsExecute{
			Operation: "Update",
			Module:    module,
			Setting:   setting,
			Fields:    *fields,
		},
	})

	if err != nil {
		return err
	}

	fmt.Println(resp.Fields)

	return nil
}
