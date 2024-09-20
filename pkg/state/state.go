package state

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"

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
	UserSessions []*types.CreateUserSessionResponse

	// The current session index
	CurrentSessionIndex int
}

// Add a new session, returns an error if token is not set
func (s *StateSessionAuth) AddSession(sess *types.CreateUserSessionResponse) error {
	s.UserSessions = append(s.UserSessions, sess)

	return nil
}

// Returns the current session
func (s *StateSessionAuth) GetCurrentSession() (*types.CreateUserSessionResponse, error) {
	if len(s.UserSessions) > s.CurrentSessionIndex-1 {
		return nil, ErrSessionNotFound
	}

	return s.UserSessions[s.CurrentSessionIndex], nil
}

// Set the current session by index
func (s *StateSessionAuth) SetCurrentSession(i int) error {
	if len(s.UserSessions) > i-1 {
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
	Persist                  *string
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

func (s *State) PersistToDisk() error {
	// Open file
	if s.Prefs.Persist == nil {
		return nil
	}

	// Get root directory from s.Prefs.Persist
	parent := filepath.Dir(*s.Prefs.Persist)
	path := filepath.Join(parent, ".evil-befall.swp")

	tmpFile, err := os.Create(path)

	if err != nil {
		return fmt.Errorf("failed to create temp file: %w", err)
	}

	defer tmpFile.Close()

	// Write to file
	if err := json.NewEncoder(tmpFile).Encode(s); err != nil {
		return fmt.Errorf("failed to write state to file: %w", err)
	}

	// Close file
	if err := tmpFile.Close(); err != nil {
		return fmt.Errorf("failed to close file: %w", err)
	}

	// Move file to final location
	if err := os.Rename(path, *s.Prefs.Persist); err != nil {
		return fmt.Errorf("failed to move file to final location: %w", err)
	}

	return nil
}

func CreateStateFromPersist(userPrefs UserPref) (*State, error) {
	if userPrefs.Persist == nil {
		return nil, fs.ErrNotExist
	}

	f, err := os.Open(*userPrefs.Persist)

	if err != nil {
		return nil, fmt.Errorf("failed to read persisted state: %w", err)
	}

	var s *State
	if err := json.NewDecoder(f).Decode(&s); err != nil {
		return nil, fmt.Errorf("failed to decode persisted state: %w", err)
	}

	return s, nil
}

func NewState(userPrefs UserPref) (*State, error) {
	if userPrefs.Persist != nil {
		s, err := CreateStateFromPersist(userPrefs)

		if err != nil && !errors.Is(err, fs.ErrNotExist) {
			return nil, fmt.Errorf("failed to create state from persisted state: %w", err)
		} else if err == nil {
			s.Prefs = userPrefs // Set user prefs to the user prefs passed in
			return s, nil
		}
	}

	return &State{
		CurrentLoc: &loc.LocMetadata{
			ID: "root",
		},
		BindAddr: "http://localhost:5173",
		Prefs:    userPrefs,
	}, nil
}
