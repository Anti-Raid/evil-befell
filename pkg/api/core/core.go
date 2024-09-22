package core

import (
	"context"

	"github.com/anti-raid/evil-befall/pkg/api"
	"github.com/anti-raid/evil-befall/pkg/fetch"
	"github.com/anti-raid/evil-befall/pkg/state"
	"github.com/anti-raid/evil-befall/types"
	"github.com/anti-raid/evil-befall/types/mewld"
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

func GetClustersHealth(ctx context.Context, state *state.State) (*mewld.InstanceList, error) {
	resp, err := fetch.Fetch(ctx, &state.StateFetchOptions, fetch.DefaultFetchOptions, fetch.FetchOptions{
		Method: "GET",
		URL:    state.StateFetchOptions.InstanceAPIUrl + "/clusters/health",
	})

	if err != nil {
		return nil, err
	}

	var instanceList mewld.InstanceList

	if err := resp.Json(&instanceList); err != nil {
		return nil, err
	}

	return &instanceList, nil
}

type GetClusterModulesData struct {
	ClusterID string `json:"path:clusterId"`
}

func GetClusterModules(ctx context.Context, state *state.State, data *GetClusterModulesData) (*[]*silverpelt.CanonicalModule, error) {
	resp, err := fetch.Fetch(ctx, &state.StateFetchOptions, fetch.DefaultFetchOptions, fetch.FetchOptions{
		Method: "GET",
		URL:    state.StateFetchOptions.InstanceAPIUrl + "/clusters/" + data.ClusterID + "/modules",
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
			api.CreateTestableRouteWithOnlyResp("getClustersHealth", GetClustersHealth),
			api.CreateTestableRouteWithReqAndResp("getClusterModules", GetClusterModules),
		),
	)
}
