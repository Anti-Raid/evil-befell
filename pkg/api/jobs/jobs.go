package jobs

import (
	"context"

	"github.com/anti-raid/evil-befall/pkg/api"
	"github.com/anti-raid/evil-befall/pkg/fetch"
	"github.com/anti-raid/evil-befall/pkg/state"
	"github.com/anti-raid/evil-befall/types"
)

type GetGuildJobData struct {
	GuildID string `json:"path:guildId"`
	JobID   string `json:"path:id"`
}

func GetGuildJob(ctx context.Context, state *state.State, data *GetGuildJobData) (*types.Job, error) {
	resp, err := fetch.Fetch(ctx, &state.StateFetchOptions, fetch.DefaultAuthorizedFetchOptions(state), fetch.FetchOptions{
		Method: "GET",
		URL:    state.StateFetchOptions.InstanceAPIUrl + "/guilds/" + data.GuildID + "/jobs/" + data.JobID,
	})

	if err != nil {
		return nil, err
	}

	var job types.Job

	if err := resp.Json(&job); err != nil {
		return nil, err
	}

	return &job, nil
}

type GetJobListData struct {
	GuildID              string `json:"path:guildId"`
	ErrorIfNoPermissions bool   `json:"query:error_if_no_permissions"`
	ErrorOnUnknownJob    bool   `json:"query:error_on_unknown_job"`
}

func GetJobList(ctx context.Context, state *state.State, data *GetJobListData) (*types.JobListResponse, error) {
	resp, err := fetch.Fetch(ctx, &state.StateFetchOptions, fetch.DefaultAuthorizedFetchOptions(state), fetch.FetchOptions{
		Method: "GET",
		URL:    state.StateFetchOptions.InstanceAPIUrl + "/guilds/" + data.GuildID + "/jobs" + api.StructToQueryParamsString(data),
	})

	if err != nil {
		return nil, err
	}

	var jobList types.JobListResponse

	if err := resp.Json(&jobList); err != nil {
		return nil, err
	}

	return &jobList, nil
}

type CreateGuildJobData struct {
	GuildID string `json:"path:guildId"`
	Name    string `json:"path:name"`
	Data    any    `json:"body:data"`
}

func CreateGuildJob(ctx context.Context, state *state.State, data *CreateGuildJobData) (*types.JobCreateResponse, error) {
	body, err := fetch.JsonBody(data.Data)

	if err != nil {
		return nil, err
	}

	resp, err := fetch.Fetch(ctx, &state.StateFetchOptions, fetch.DefaultAuthorizedFetchOptions(state), fetch.FetchOptions{
		Method: "POST",
		URL:    state.StateFetchOptions.InstanceAPIUrl + "/guilds/" + data.GuildID + "/jobs/" + data.Name,
		Body:   body,
	})

	if err != nil {
		return nil, err
	}

	var jcr types.JobCreateResponse

	if err := resp.Json(&jcr); err != nil {
		return nil, err
	}

	return &jcr, nil
}

// GetIOAuthDownloadLinkData needs to return a struct as it is special
type GetIOAuthDownloadLinkData struct {
	ID         string `json:"path:id"`
	NoRedirect bool   `json:"query:no_redirect"`
}

func GetIOAuthDownloadLink(ctx context.Context, state *state.State, data *GetIOAuthDownloadLinkData) (*string, error) {
	url := state.StateFetchOptions.InstanceAPIUrl + "/jobs/" + data.ID + "/ioauth/download-link" + api.StructToQueryParamsString(data)
	return &url, nil
}

func init() {
	api.RegisterTestableRouteCategory(
		api.NewTestableRouteCategory(
			"jobs",
			api.CreateTestableRouteWithReqAndResp("getGuildJob", GetGuildJob),
			api.CreateTestableRouteWithReqAndResp("getJobList", GetJobList),
			api.CreateTestableRouteWithReqAndResp("createGuildJob", CreateGuildJob),
			api.CreateTestableRouteWithReqAndResp("getIOAuthDownloadLink", GetIOAuthDownloadLink),
		),
	)
}
