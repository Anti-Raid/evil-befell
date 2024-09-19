package auth

import (
	"context"

	"github.com/anti-raid/evil-befall/pkg/fetch"
	"github.com/anti-raid/evil-befall/pkg/state"
	"github.com/anti-raid/evil-befall/types"
)

func CreateOauth2Login(ctx context.Context, state *state.State, data types.AuthorizeRequest) (*types.UserLogin, error) {
	body, err := fetch.JsonBody(data)

	if err != nil {
		return nil, err
	}

	resp, err := fetch.Fetch(ctx, &state.StateFetchOptions, fetch.DefaultFetchOptions, fetch.FetchOptions{
		Method: "POST",
		URL:    state.StateFetchOptions.InstanceAPIUrl + "/oauth2",
		Body:   body,
	})

	if err != nil {
		return nil, err
	}

	var userLogin types.UserLogin

	if err := resp.Json(&userLogin); err != nil {
		return nil, err
	}

	return &userLogin, nil
}
