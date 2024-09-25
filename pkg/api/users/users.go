package users

import (
	"context"

	"github.com/anti-raid/evil-befall/pkg/api"
	"github.com/anti-raid/evil-befall/pkg/fetch"
	"github.com/anti-raid/evil-befall/pkg/state"
	"github.com/anti-raid/evil-befall/types"
)

type GetUserData struct {
	ID string `json:"path:id"`
}

func GetUser(ctx context.Context, state *state.State, data *GetUserData) (*types.User, error) {
	resp, err := fetch.Fetch(ctx, &state.StateFetchOptions, fetch.DefaultFetchOptions, fetch.FetchOptions{
		Method: "GET",
		URL:    state.StateFetchOptions.InstanceAPIUrl + "/users/" + data.ID,
	})

	if err != nil {
		return nil, err
	}

	var user types.User

	if err := resp.Json(&user); err != nil {
		return nil, err
	}

	return &user, nil
}

type GetUserGuildsData struct {
	Refresh bool `json:"query:refresh"`
}

func GetUserGuilds(ctx context.Context, state *state.State, data *GetUserGuildsData) (*types.DashboardGuildData, error) {
	resp, err := fetch.Fetch(ctx, &state.StateFetchOptions, fetch.DefaultAuthorizedFetchOptions(state), fetch.FetchOptions{
		Method: "GET",
		URL:    state.StateFetchOptions.InstanceAPIUrl + "/users/@me/guilds" + api.StructToQueryParamsString(data),
	})

	if err != nil {
		return nil, err
	}

	var guilds types.DashboardGuildData

	if err := resp.Json(&guilds); err != nil {
		return nil, err
	}

	return &guilds, nil
}

type GetUserGuildBaseInfoData struct {
	GuildID string `json:"path:guildId"`
}

func GetUserGuildBaseInfo(ctx context.Context, state *state.State, data *GetUserGuildBaseInfoData) (*types.DashboardGuild, error) {
	resp, err := fetch.Fetch(ctx, &state.StateFetchOptions, fetch.DefaultAuthorizedFetchOptions(state), fetch.FetchOptions{
		Method: "GET",
		URL:    state.StateFetchOptions.InstanceAPIUrl + "/users/@me/guilds/" + data.GuildID,
	})

	if err != nil {
		return nil, err
	}

	var guild types.DashboardGuild

	if err := resp.Json(&guild); err != nil {
		return nil, err
	}

	return &guild, nil
}

func init() {
	api.RegisterTestableRouteCategory(
		api.NewTestableRouteCategory(
			"user",
			api.CreateTestableRouteWithReqAndResp("getUser", GetUser),
			api.CreateTestableRouteWithReqAndResp("getUserGuilds", GetUserGuilds),
			api.CreateTestableRouteWithReqAndResp("getUserGuildBaseInfo", GetUserGuildBaseInfo),
		),
	)
}
