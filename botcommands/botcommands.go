package botcommands

import (
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"sort"

	"github.com/gin-gonic/gin"

	"github.com/bwmarrin/discordgo"
)

// Global instance of our tracker.
var tracker = InitiativeTracker{
	players: make(map[string]*Player),
}

func handleIoCommand(s *discordgo.Session, i *discordgo.InteractionCreate, vs *discordgo.VoiceState) {
	data := i.ApplicationCommandData()
	initiativeRoll := data.Options[0].IntValue()
	user := i.Member.User

	tracker.mu.Lock()
	defer tracker.mu.Unlock()

	// If this is the first person to set initiative, they set the "active" channel.
	// Or if the tracker was reset, set the new active channel.
	if tracker.activeChannelID == "" {
		tracker.activeChannelID = vs.ChannelID
	} else if tracker.activeChannelID != vs.ChannelID {
		// If a user from another channel tries to set IO, tell them a game is active elsewhere.
		sendEphemeralResponse(s, i, "An initiative order is already active in another channel. Use `/io-reset` in that channel to clear it.")
		return
	}

	// Create or update the player's data.
	player := &Player{
		UserID:     user.ID,
		Username:   user.Username,
		AvatarURL:  user.AvatarURL("128"), // Get a 128x128 avatar URL
		Initiative: int(initiativeRoll),
	}
	tracker.players[user.ID] = player

	log.Printf("Updated initiative for %s to %d in channel %s", user.Username, initiativeRoll, vs.ChannelID)
	for _, p := range tracker.players {
		log.Printf("Player: %s, Initiative: %d", p.Username, p.Initiative)
	}

	// Respond to the user to confirm.
	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: fmt.Sprintf("Your initiative is set to **%d**!", initiativeRoll),
		},
	})
}

func handleIoShowCommand(s *discordgo.Session, i *discordgo.InteractionCreate, vs *discordgo.VoiceState) {
	log.Println("Handling /io-show command started.")
	tracker.mu.RLock()
	defer tracker.mu.RUnlock()

	// Only allow showing IO if the user is in the active channel.
	log.Printf("checking channel ID: %s, User's voice channel ID: %s", tracker.activeChannelID, vs.ChannelID)
	if tracker.activeChannelID == "" || tracker.activeChannelID != vs.ChannelID {
		sendEphemeralResponse(s, i, "You must be in the active game's voice channel to view the initiative order.")
		return
	}

	// Create a formatted string of all players and their initiatives.
	var response string
	for _, player := range tracker.players {
		response += fmt.Sprintf("%s: **%d**\n", player.Username, player.Initiative)
	}

	if response == "" {
		response = "No players have set their initiative yet."
	}

	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: response,
		},
	})
}

func handleIoResetCommand(s *discordgo.Session, i *discordgo.InteractionCreate, vs *discordgo.VoiceState) {
	tracker.mu.Lock()
	defer tracker.mu.Unlock()

	// Only allow reset from the active channel, if one is set.
	if tracker.activeChannelID != "" && tracker.activeChannelID != vs.ChannelID {
		sendEphemeralResponse(s, i, "You must be in the active game's voice channel to reset the initiative order.")
		return
	}

	// Clear the players map and the active channel ID.
	tracker.players = make(map[string]*Player)
	tracker.activeChannelID = ""

	log.Printf("Initiative order reset by %s for channel %s", i.Member.User.Username, vs.ChannelID)

	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: "The initiative order has been reset!",
		},
	})
}

func handleIoAddCommand(s *discordgo.Session, i *discordgo.InteractionCreate, vs *discordgo.VoiceState) {
	data := i.ApplicationCommandData()
	npcName := data.Options[0].StringValue()

	tracker.mu.Lock()
	defer tracker.mu.Unlock()

	// Only allow adding NPCs if the user is in the active channel.
	if tracker.activeChannelID == "" || tracker.activeChannelID != vs.ChannelID {
		sendEphemeralResponse(s, i, "You must be in the active game's voice channel to add an NPC.")
		return
	}

	// Create an npc player with a unique ID. and random initiative.
	npcPlayer := &Player{
		UserID:     fmt.Sprintf("npc-%s", npcName), // Unique ID for NPCs
		Username:   npcName,
		AvatarURL:  "",                // No avatar for NPCs
		Initiative: rand.Intn(23) + 1, // Random initiative between 1 and 23
	}
	tracker.players[npcPlayer.UserID] = npcPlayer

	log.Printf("Added NPC %s with initiative %d", npcName, npcPlayer.Initiative)

	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: fmt.Sprintf("NPC **%s** added with initiative **%d**!", npcName, npcPlayer.Initiative),
		},
	})
}

// Helper to find a user's voice state.
func findUserVoiceState(s *discordgo.Session, guildID, userID string) (*discordgo.VoiceState, error) {
	guild, err := s.State.Guild(guildID)
	if err != nil {
		return nil, err
	}
	for _, vs := range guild.VoiceStates {
		if vs.UserID == userID {
			return vs, nil
		}
	}
	return nil, fmt.Errorf("could not find user %s in a voice channel", userID)
}

// Helper to send a private message back to the user.
func sendEphemeralResponse(s *discordgo.Session, i *discordgo.InteractionCreate, content string) {
	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: content,
			Flags:   discordgo.MessageFlagsEphemeral,
		},
	})
}

func InitiativeHandler(c *gin.Context) {
	log.Print("InitiativeHandler called")

	// Set headers for CORS and JSON content type.
	c.Header("Access-Control-Allow-Origin", "*")
	c.Header("Content-Type", "application/json")

	tracker.mu.RLock()
	defer tracker.mu.RUnlock()

	guildID := Getenvvar("XGuildID")
	if guildID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "XGuildID environment variable not set"})
		return
	}

	// We only want to show players from the currently active channel.
	var activePlayers []*Player
	for _, p := range tracker.players {
		// For NPCs (which have UserIDs starting with "npc-"), add them directly
		if len(p.UserID) >= 4 && p.UserID[:4] == "npc-" {
			activePlayers = append(activePlayers, p)
			continue
		}
		// For regular players, check if they're still in the active channel
		vs, err := findUserVoiceState(DiscordSession, guildID, p.UserID)
		if err == nil && vs.ChannelID == tracker.activeChannelID {
			activePlayers = append(activePlayers, p)
		}
	}

	// If the active channel is set, find all users in it to construct the list.
	// This ensures we show everyone in the channel, even if they haven't set initiative yet.
	guild, err := DiscordSession.State.Guild(guildID)
	// If we can't find the guild, return an error.
	if err == nil && tracker.activeChannelID != "" {
		allInChannel := make(map[string]bool)
		for _, vs := range guild.VoiceStates {
			if vs.ChannelID == tracker.activeChannelID {
				allInChannel[vs.UserID] = true
				// If this user isn't in our tracker yet, add them with 0 initiative.
				if _, ok := tracker.players[vs.UserID]; !ok {
					user, _ := DiscordSession.User(vs.UserID)
					// If we can't find the user, skip adding them.
					if user == nil {
						continue
					}
					// Add the user with an initiative of 0.
					tracker.players[vs.UserID] = &Player{
						UserID:     user.ID,
						Username:   user.Username,
						AvatarURL:  user.AvatarURL("128"),
						Initiative: 0, // Default initiative
					}
				}
			}
		}
		// Filter out players who are no longer in the channel, but keep NPCs
		currentPlayers := make([]*Player, 0)
		for userID, player := range tracker.players {
			// Keep NPCs (which have UserIDs starting with "npc-")
			if len(userID) >= 4 && userID[:4] == "npc-" {
				currentPlayers = append(currentPlayers, player)
			} else if _, inChannel := allInChannel[userID]; inChannel {
				currentPlayers = append(currentPlayers, player)
			}
		}
		activePlayers = currentPlayers
	}

	// Sort players by initiative, descending.
	sort.Slice(activePlayers, func(i, j int) bool {
		return activePlayers[i].Initiative > activePlayers[j].Initiative
	})

	// Encode the sorted list to JSON and send it.
	c.JSON(http.StatusOK, activePlayers)
}
