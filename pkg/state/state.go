package state

import (
	"errors"

	"github.com/anti-raid/evil-befall/pkg/loc"
	"github.com/anti-raid/evil-befall/types"
)

var (
	ErrSessionHasNoToken = errors.New("session does not have a token")
	ErrSessionNotFound   = errors.New("session was not found")
)

type StateSessionAuth struct {
	// The session information for the logged in user
	//
	// Because Evil Befall can technically support multiple sessions, we can just store the raw UserSessionList
	// object here
	UserSessions *types.UserSessionList

	// The current session index
	CurrentSessionIndex int
}

// Add a new session, returns an error if token is not set
func (s *StateSessionAuth) AddSession(sess *types.UserSession) error {
	if sess.Token == nil {
		return ErrSessionHasNoToken
	}

	s.UserSessions.Sessions = append(s.UserSessions.Sessions, sess)

	return nil
}

// Returns the current session
func (s *StateSessionAuth) GetCurrentSession() (*types.UserSession, error) {
	if len(s.UserSessions.Sessions) > s.CurrentSessionIndex-1 {
		return nil, ErrSessionNotFound
	}

	return s.UserSessions.Sessions[s.CurrentSessionIndex], nil
}

// Set the current session by index
func (s *StateSessionAuth) SetCurrentSession(i int) error {
	if len(s.UserSessions.Sessions) > i-1 {
		return ErrSessionNotFound
	}

	s.CurrentSessionIndex = i

	return nil
}

// Returns if the user is currently authorized into a session
func (s *StateSessionAuth) IsAuthorized() bool {
	_, err := s.GetCurrentSession()

	return !errors.Is(err, ErrSessionNotFound)
}

type StateFetchOptions struct {
	// The API URL for the Anti-Raid instance
	InstanceAPIUrl string
}

type UserPref struct {
	MouseEnabledInTView      bool
	PasteEnabledInTView      bool
	FullscreenEnabledInTView bool
}

// Stores all the state for the application
type State struct {
	// The current location Evil Befall is at
	CurrentLoc *loc.LocMetadata

	// Session auth
	Session StateSessionAuth

	// State fetch options
	StateFetchOptions StateFetchOptions

	BindAddr string // Bind address with login/logout etc.

	Prefs UserPref
}

func NewState() *State {
	return &State{
		CurrentLoc: &loc.LocMetadata{
			ID: "root",
		},
		BindAddr: "http://localhost:5173",
	}
}
