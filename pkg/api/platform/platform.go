package platform

import (
	"context"

	"github.com/anti-raid/evil-befall/pkg/api"
	"github.com/anti-raid/evil-befall/pkg/fetch"
	"github.com/anti-raid/evil-befall/pkg/state"
	"github.com/anti-raid/evil-befall/types/dovetypes"
)

type GetPlatformUserData struct {
	ID       string `json:"path:id"`
	Platform string `json:"query:platform"`
}

func GetPlatformUser(ctx context.Context, state *state.State, data *GetPlatformUserData) (*dovetypes.PlatformUser, error) {
	resp, err := fetch.Fetch(ctx, &state.StateFetchOptions, fetch.DefaultFetchOptions, fetch.FetchOptions{
		Method: "GET",
		URL:    state.StateFetchOptions.InstanceAPIUrl + "/platform/user/" + data.ID + api.StructToQueryParamsString(data),
	})

	if err != nil {
		return nil, err
	}

	var platformUser dovetypes.PlatformUser

	if err := resp.Json(&platformUser); err != nil {
		return nil, err
	}

	return &platformUser, nil
}

type ClearPlatformUserCacheData struct {
	ID       string `json:"path:id"`
	Platform string `json:"query:platform"`
}

func ClearPlatformUserCache(ctx context.Context, state *state.State, data *ClearPlatformUserCacheData) error {
	_, err := fetch.Fetch(ctx, &state.StateFetchOptions, fetch.DefaultFetchOptions, fetch.FetchOptions{
		Method: "DELETE",
		URL:    state.StateFetchOptions.InstanceAPIUrl + "/platform/user/" + data.ID + api.StructToQueryParamsString(data),
	})

	if err != nil {
		return err
	}

	return nil
}

func init() {
	api.RegisterTestableRouteCategory(
		api.NewTestableRouteCategory(
			"platform",
			api.CreateTestableRouteWithReqAndResp("getPlatformUser", GetPlatformUser),
			api.CreateTestableRouteWithOnlyReq("clearPlatformUserCache", ClearPlatformUserCache),
		),
	)
}
