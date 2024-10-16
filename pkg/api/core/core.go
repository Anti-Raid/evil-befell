package core

import (
	"context"

	"github.com/anti-raid/evil-befall/pkg/api"
	"github.com/anti-raid/evil-befall/pkg/fetch"
	"github.com/anti-raid/evil-befall/pkg/state"
	"github.com/anti-raid/evil-befall/types"
	"github.com/anti-raid/evil-befall/types/silverpelt"
)

func GetApiConfig(ctx context.Context, state *state.State) (*types.ApiConfig, error) {
	resp, err := fetch.Fetch(ctx, &state.StateFetchOptions, fetch.DefaultFetchOptions, fetch.FetchOptions{
		Method: "GET",
		URL:    state.StateFetchOptions.InstanceAPIUrl + "/config",
	})

	if err != nil {
		return nil, err
	}

	var apiConfig types.ApiConfig

	if err := resp.Json(&apiConfig); err != nil {
		return nil, err
	}

	return &apiConfig, nil
}

func GetModules(ctx context.Context, state *state.State) (*[]*silverpelt.CanonicalModule, error) {
	resp, err := fetch.Fetch(ctx, &state.StateFetchOptions, fetch.DefaultFetchOptions, fetch.FetchOptions{
		Method: "GET",
		URL:    state.StateFetchOptions.InstanceAPIUrl + "/modules",
	})

	if err != nil {
		return nil, err
	}

	var moduleList []*silverpelt.CanonicalModule

	if err := resp.Json(&moduleList); err != nil {
		return nil, err
	}

	return &moduleList, nil
}

func init() {
	api.RegisterTestableRouteCategory(
		api.NewTestableRouteCategory(
			"core",
			api.CreateTestableRouteWithOnlyResp("getApiConfig", GetApiConfig),
			api.CreateTestableRouteWithOnlyResp("getModules", GetModules),
		),
	)
}
