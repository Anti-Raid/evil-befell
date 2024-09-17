package state

import (
	"encoding/json"
	"errors"

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

type StateFetchOptions struct {
	// The API URL for the Anti-Raid instance
	InstanceAPIUrl string
}

// Stores all the state for the application
type State struct {
	// The current location Evil Befall is at
	CurrentLoc string

	// Session auth
	Session StateSessionAuth

	// Per location JSON data
	LocationData map[string]json.RawMessage

	// State fetch options
	StateFetchOptions StateFetchOptions
}

// Add a new session, returns an error if token is not set
func (s *State) AddSession(sess *types.UserSession) error {
	if sess.Token == nil {
		return ErrSessionHasNoToken
	}

	s.Session.UserSessions.Sessions = append(s.Session.UserSessions.Sessions, sess)

	return nil
}

// Set the current session by index
func (s *State) SetCurrentSession(i int) error {
	if len(s.Session.UserSessions.Sessions) > i-1 {
		return ErrSessionNotFound
	}

	s.Session.CurrentSessionIndex = i

	return nil
}

// Returns the current session
func (s *State) GetCurrentSession() (*types.UserSession, error) {
	if len(s.Session.UserSessions.Sessions) > s.Session.CurrentSessionIndex-1 {
		return nil, ErrSessionNotFound
	}

	return s.Session.UserSessions.Sessions[s.Session.CurrentSessionIndex], nil
}
