package auth

import (
	"context"
	"errors"
	"net/http"
	"strings"

	"github.com/anti-raid/evil-befall/pkg/api/auth"
	"github.com/anti-raid/evil-befall/pkg/state"
	"github.com/anti-raid/evil-befall/types"
)

const NoBotInvite = "https://discord.com/api/oauth2/authorize?client_id={client_id}&response_type=code&redirect_uri={redirect_url}&scope=guilds+identify&prompt=none"

func GetAuthURL(ctx context.Context, state *state.State, apiConfig *types.ApiConfig) string {
	inviteURL := NoBotInvite

	inviteURL = strings.Replace(inviteURL, "{client_id}", apiConfig.ClientID, 1)
	inviteURL = strings.Replace(inviteURL, "{redirect_url}", state.BindAddr+"/authorize", 1)

	return inviteURL
}

func CreateSessionOnServerAddr(ctx context.Context, state *state.State) (*types.CreateUserSessionResponse, error) {
	var createdSessionChan = make(chan *types.CreateUserSessionResponse)

	mux := http.NewServeMux()

	bindAddr := strings.Split(state.BindAddr, ":")

	server := &http.Server{
		Addr:    ":" + bindAddr[len(bindAddr)-1],
		Handler: mux,
	}

	mux.HandleFunc("/authorize", func(w http.ResponseWriter, r *http.Request) {
		code := r.URL.Query().Get("code")

		createdSession, err := auth.CreateOauth2Login(ctx, state, types.AuthorizeRequest{
			Code:        code,
			RedirectURI: state.BindAddr + "/authorize",
			Protocol:    "a1",
			Scope:       "normal",
		})

		if err != nil {
			w.Write([]byte("Error: " + err.Error()))
			return
		}

		w.Write([]byte("Success! You can close this window now."))

		createdSessionChan <- createdSession
	})

	go func() {
		defer recover()

		if err := server.ListenAndServe(); !errors.Is(err, http.ErrServerClosed) {
			panic(err)
		}
	}()

	createdSession := <-createdSessionChan

	if err := server.Shutdown(ctx); err != nil {
		return nil, err
	}

	return createdSession, nil
}
