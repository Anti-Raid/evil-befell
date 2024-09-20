package showstate

import (
	"encoding/json"
	"fmt"

	"github.com/anti-raid/evil-befall/pkg/state"
)

type ShowStateRoute struct {
}

func (r *ShowStateRoute) Command() string {
	return "showstate"
}

func (r *ShowStateRoute) Description() string {
	return "Prints out the current state"
}

func (r *ShowStateRoute) Arguments() [][3]string {
	return [][3]string{}
}

func (r *ShowStateRoute) Setup(state *state.State) error {
	return nil
}

func (r *ShowStateRoute) Destroy(state *state.State) error {
	return nil
}

func (r *ShowStateRoute) Render(state *state.State, args map[string]string) error {
	indentedJson, err := json.MarshalIndent(state, "", "  ")

	if err != nil {
		return err
	}

	fmt.Println(string(indentedJson))

	return nil
}
