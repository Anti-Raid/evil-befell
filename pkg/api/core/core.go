package core

import (
	"context"

	"github.com/anti-raid/evil-befall/pkg/fetch"
	"github.com/anti-raid/evil-befall/pkg/state"
	"github.com/anti-raid/evil-befall/types"
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
