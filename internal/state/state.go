package state

import "github.com/anti-raid/evil-befall/types"

// Stores all the state for the application
type State struct {
	// The current location Evil Befall is at
	CurrentLoc string

	// The session information for the logged in user
	//
	// Because Evil Befall can technically support multiple sessions, we can just store the raw UserSessionList
	// object here
	UserSessions *types.UserSessionList

	// The current session index
	CurrentSessionIndex int
}
