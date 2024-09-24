package guilds

import (
	"context"

	"github.com/anti-raid/evil-befall/pkg/api"
	"github.com/anti-raid/evil-befall/pkg/fetch"
	"github.com/anti-raid/evil-befall/pkg/state"
	"github.com/anti-raid/evil-befall/types"
	"github.com/anti-raid/evil-befall/types/silverpelt"
)

type GetStaffTeamData struct {
	GuildID string `json:"path:guildId"`
}

func GetStaffTeam(ctx context.Context, state *state.State, data *GetStaffTeamData) (*types.GuildStaffTeam, error) {
	resp, err := fetch.Fetch(ctx, &state.StateFetchOptions, fetch.DefaultFetchOptions, fetch.FetchOptions{
		Method: "GET",
		URL:    state.StateFetchOptions.InstanceAPIUrl + "/guilds/" + data.GuildID + "/staff-team",
	})

	if err != nil {
		return nil, err
	}

	var staffTeam types.GuildStaffTeam

	if err := resp.Json(&staffTeam); err != nil {
		return nil, err
	}

	return &staffTeam, nil
}

type GetModuleConfigurationsData struct {
	GuildID string `json:"path:guildId"`
}

func GetModuleConfigurations(ctx context.Context, state *state.State, data *GetModuleConfigurationsData) (*[]*silverpelt.GuildModuleConfiguration, error) {
	resp, err := fetch.Fetch(ctx, &state.StateFetchOptions, fetch.DefaultAuthorizedFetchOptions(state), fetch.FetchOptions{
		Method: "GET",
		URL:    state.StateFetchOptions.InstanceAPIUrl + "/guilds/" + data.GuildID + "/module-configurations",
	})

	if err != nil {
		return nil, err
	}

	var gmc []*silverpelt.GuildModuleConfiguration

	if err := resp.Json(&gmc); err != nil {
		return nil, err
	}

	return &gmc, nil
}

type PatchModuleConfigurationsData struct {
	*types.PatchGuildModuleConfiguration `json:"patch"`
	GuildID                              string `json:"path:guildId"`
}

func PatchModuleConfiguration(ctx context.Context, state *state.State, data *PatchModuleConfigurationsData) (*silverpelt.GuildModuleConfiguration, error) {
	body, err := fetch.JsonBody(data.PatchGuildModuleConfiguration)

	if err != nil {
		return nil, err
	}

	resp, err := fetch.Fetch(ctx, &state.StateFetchOptions, fetch.DefaultAuthorizedFetchOptions(state), fetch.FetchOptions{
		Method: "PATCH",
		URL:    state.StateFetchOptions.InstanceAPIUrl + "/guilds/" + data.GuildID + "/module-configurations",
		Body:   body,
	})

	if err != nil {
		return nil, err
	}

	var gmc silverpelt.GuildModuleConfiguration

	if err := resp.Json(&gmc); err != nil {
		return nil, err
	}

	return &gmc, nil
}

type GetAllCommandConfigurationsData struct {
	GuildID string `json:"path:guildId"`
}

func GetAllCommandConfigurations(ctx context.Context, state *state.State, data *GetAllCommandConfigurationsData) (*[]*silverpelt.FullGuildCommandConfiguration, error) {
	resp, err := fetch.Fetch(ctx, &state.StateFetchOptions, fetch.DefaultAuthorizedFetchOptions(state), fetch.FetchOptions{
		Method: "GET",
		URL:    state.StateFetchOptions.InstanceAPIUrl + "/guilds/" + data.GuildID + "/command-configurations",
	})

	if err != nil {
		return nil, err
	}

	var gcc []*silverpelt.FullGuildCommandConfiguration

	if err := resp.Json(&gcc); err != nil {
		return nil, err
	}

	return &gcc, nil
}

type PatchCommandConfigurationsData struct {
	*types.PatchGuildCommandConfiguration `json:"patch"`
	GuildID                               string `json:"path:guildId"`
}

func PatchCommandConfiguration(ctx context.Context, state *state.State, data *PatchCommandConfigurationsData) (*silverpelt.FullGuildCommandConfiguration, error) {
	body, err := fetch.JsonBody(data.PatchGuildCommandConfiguration)

	if err != nil {
		return nil, err
	}

	resp, err := fetch.Fetch(ctx, &state.StateFetchOptions, fetch.DefaultAuthorizedFetchOptions(state), fetch.FetchOptions{
		Method: "PATCH",
		URL:    state.StateFetchOptions.InstanceAPIUrl + "/guilds/" + data.GuildID + "/command-configurations",
		Body:   body,
	})

	if err != nil {
		return nil, err
	}

	var gcc silverpelt.FullGuildCommandConfiguration

	if err := resp.Json(&gcc); err != nil {
		return nil, err
	}

	return &gcc, nil
}

type SettingsExecuteData struct {
	SettingsExecuteData *types.SettingsExecute `json:"body"`
	GuildID             string                 `json:"path:guildId"`
}

func SettingsExecute(ctx context.Context, state *state.State, data *SettingsExecuteData) (*types.SettingsExecuteResponse, error) {
	body, err := fetch.JsonBody(data.SettingsExecuteData)

	if err != nil {
		return nil, err
	}

	resp, err := fetch.Fetch(ctx, &state.StateFetchOptions, fetch.DefaultAuthorizedFetchOptions(state), fetch.FetchOptions{
		Method: "POST",
		URL:    state.StateFetchOptions.InstanceAPIUrl + "/guilds/" + data.GuildID + "/settings",
		Body:   body,
	})

	if err != nil {
		return nil, err
	}

	var ser types.SettingsExecuteResponse

	if err := resp.Json(&ser); err != nil {
		return nil, err
	}

	return &ser, nil
}

func init() {
	api.RegisterTestableRouteCategory(
		api.NewTestableRouteCategory(
			"guilds",
			api.CreateTestableRouteWithReqAndResp("getStaffTeam", GetStaffTeam),
			api.CreateTestableRouteWithReqAndResp("getModuleConfigurations", GetModuleConfigurations),
			api.CreateTestableRouteWithReqAndResp("patchModuleConfiguration", PatchModuleConfiguration),
			api.CreateTestableRouteWithReqAndResp("getAllCommandConfigurations", GetAllCommandConfigurations),
			api.CreateTestableRouteWithReqAndResp("patchCommandConfiguration", PatchCommandConfiguration),
			api.CreateTestableRouteWithReqAndResp("settingsExecute", SettingsExecute),
		),
	)
}
