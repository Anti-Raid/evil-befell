package auth

import (
	"context"

	"github.com/anti-raid/evil-befall/pkg/api"
	"github.com/anti-raid/evil-befall/pkg/fetch"
	"github.com/anti-raid/evil-befall/pkg/state"
	"github.com/anti-raid/evil-befall/types"
)

// CreateIoAuthLogin is special because it primarily uses query parameters
type CreateIoAuthLoginData struct {
	PathRedirectData string  `json:"query:path_rd"`
	PathCode         *string `json:"query:path_code"`
}

func CreateIoAuthLogin(ctx context.Context, state *state.State, data *CreateIoAuthLoginData) (*string, error) {
	url := state.StateFetchOptions.InstanceAPIUrl + "/ioauth/login" + api.StructToQueryParamsString(data)
	return &url, nil
}

func TestAuth(ctx context.Context, state *state.State, data *types.TestAuth) (*types.TestAuthResponse, error) {
	body, err := fetch.JsonBody(data)

	if err != nil {
		return nil, err
	}

	resp, err := fetch.Fetch(ctx, &state.StateFetchOptions, fetch.DefaultFetchOptions, fetch.FetchOptions{
		Method: "POST",
		URL:    state.StateFetchOptions.InstanceAPIUrl + "/auth/test",
		Body:   body,
	})

	if err != nil {
		return nil, err
	}

	var res types.TestAuthResponse

	if err := resp.Json(&res); err != nil {
		return nil, err
	}

	return &res, nil
}

func CreateOauth2Login(ctx context.Context, state *state.State, data types.AuthorizeRequest) (*types.CreateUserSessionResponse, error) {
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

	var res types.CreateUserSessionResponse

	if err := resp.Json(&res); err != nil {
		return nil, err
	}

	return &res, nil
}

func GetUserSessions(ctx context.Context, state *state.State) (*types.UserSessionList, error) {
	resp, err := fetch.Fetch(ctx, &state.StateFetchOptions, fetch.DefaultAuthorizedFetchOptions(state), fetch.FetchOptions{
		Method: "GET",
		URL:    state.StateFetchOptions.InstanceAPIUrl + "/sessions",
	})

	if err != nil {
		return nil, err
	}

	var res types.UserSessionList

	if err := resp.Json(&res); err != nil {
		return nil, err
	}

	return &res, nil
}

func CreateUserSession(ctx context.Context, state *state.State, data *types.CreateUserSession) (*types.CreateUserSessionResponse, error) {
	body, err := fetch.JsonBody(data)

	if err != nil {
		return nil, err
	}

	resp, err := fetch.Fetch(ctx, &state.StateFetchOptions, fetch.DefaultAuthorizedFetchOptions(state), fetch.FetchOptions{
		Method: "POST",
		URL:    state.StateFetchOptions.InstanceAPIUrl + "/sessions",
		Body:   body,
	})

	if err != nil {
		return nil, err
	}

	var res types.CreateUserSessionResponse

	if err := resp.Json(&res); err != nil {
		return nil, err
	}

	return &res, nil
}

type RevokeUserSessionData struct {
	SessionID string `json:"path:session_id"`
}

func RevokeUserSession(ctx context.Context, state *state.State, data *RevokeUserSessionData) error {
	_, err := fetch.Fetch(ctx, &state.StateFetchOptions, fetch.DefaultAuthorizedFetchOptions(state), fetch.FetchOptions{
		Method: "DELETE",
		URL:    state.StateFetchOptions.InstanceAPIUrl + "/sessions/" + data.SessionID,
	})

	if err != nil {
		return err
	}

	state.Session.RemoveSessionIfExists(data.SessionID)

	return nil
}

func init() {
	api.RegisterTestableRouteCategory(
		api.NewTestableRouteCategory(
			"auth",
			api.CreateTestableRouteWithReqAndResp("createIoAuthLogin", CreateIoAuthLogin),
			api.CreateTestableRouteWithReqAndResp("testAuth", TestAuth),
			api.CreateTestableRouteWithReqAndResp("createOauth2Login", CreateOauth2Login),
			api.CreateTestableRouteWithOnlyResp("getUserSessions", GetUserSessions),
			api.CreateTestableRouteWithReqAndResp("createUserSession", CreateUserSession),
			api.CreateTestableRouteWithOnlyReq("revokeUserSession", RevokeUserSession),
		),
	)
}
