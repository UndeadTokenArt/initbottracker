package botcommands

import (
	"sync"
)

// Player holds the data for a single player in the initiative order.
type Player struct {
	UserID     string `json:"userID"`
	Username   string `json:"username"`
	AvatarURL  string `json:"avatarURL"`
	Initiative int    `json:"initiative"`
}

// InitiativeTracker safely manages the state of our players.
type InitiativeTracker struct {
	mu      sync.RWMutex
	players map[string]*Player
	// We store the voice channel ID to keep track of the "active" game.
	// A more complex bot could support multiple channels at once.
	activeChannelID string
}
