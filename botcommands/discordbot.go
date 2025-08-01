package botcommands

import (
	"log"
	"os"

	"github.com/bwmarrin/discordgo"
	"github.com/joho/godotenv"
)

func Getenvvar(envVar string) string {
	// load the bot token from the environment variable.
	err := godotenv.Load()
	dir, _ := os.Getwd()
	log.Printf("Current working directory: %s", dir)
	if _, err := os.Stat(".env"); err == nil {
		log.Println(".env file found")
	} else {
		log.Println(".env file NOT found")
	}

	if err != nil {
		log.Println("Note: error loading .env file, using environment variables directly.")
	} else {
		log.Println("Loaded environment variables from .env file.")
	}

	token := os.Getenv(envVar)
	if token == "" {
		log.Fatalf("%v environment variable not set.", envVar)
	} else {
		log.Printf("%v environment variable loaded successfully.", envVar)
	}
	return token
}

// global instance of Discord session
var DiscordSession *discordgo.Session

// createDiscordSession initializes a new Discord session and registers the necessary event handlers.
func CreateDiscordSession() (*discordgo.Session, error) {
	token := Getenvvar("DISCORD_TOKEN")

	// Create a new Discord session using the provided bot token.
	dg, err := discordgo.New("Bot " + token)
	if err != nil {
		log.Fatalf("Error creating Discord session: %v", err)
	} else {
		log.Println("Discord session created successfully.")
	}
	// Store the global Discord session.
	DiscordSession = dg

	// Register the interactionCreate func as a callback for InteractionCreate events.
	dg.AddHandler(interactionCreate)

	// Register the ready func as a callback for the Ready event.
	// This is where we will register our slash commands.
	dg.AddHandler(ready)

	// We need intents to see guild members and their voice states.
	dg.Identify.Intents = discordgo.IntentsGuilds | discordgo.IntentsGuildVoiceStates | discordgo.IntentsGuildMembers

	// Open a websocket connection to Discord and begin listening.
	err = dg.Open()
	if err != nil {
		log.Fatalf("Error opening connection: %v", err)
	}
	defer dg.Close()
	return dg, nil
}

// This function is called when the bot is ready to start working.
func ready(s *discordgo.Session, event *discordgo.Ready) {
	log.Println("Bot is ready. Registering commands...")

	// Define the /io slash command.
	cmd := &discordgo.ApplicationCommand{
		Name:        "io",
		Description: "Set your initiative order.",
		Options: []*discordgo.ApplicationCommandOption{
			{
				Type:        discordgo.ApplicationCommandOptionInteger,
				Name:        "roll",
				Description: "Your initiative roll (e.g., 14)",
				Required:    true,
			},
		},
	}

	// Define the /io-reset slash command.
	resetCmd := &discordgo.ApplicationCommand{
		Name:        "io-reset",
		Description: "Clears the current initiative order for your channel.",
	}

	showCmd := &discordgo.ApplicationCommand{
		Type:        discordgo.ChatApplicationCommand,
		Name:        "ioshow",
		Description: "Shows the current initiative order for your channel.",
	}

	// Register the commands. For a personal bot, you can register it for a specific
	// guild for instant updates. Global registration can take up to an hour.
	// _, err := s.ApplicationCommandCreate(s.State.User.ID, "YOUR_GUILD_ID_HERE", cmd)
	_, err := s.ApplicationCommandCreate(s.State.User.ID, "", cmd)
	if err != nil {
		log.Printf("Cannot create '/io' command: %v", err)
	}

	_, err = s.ApplicationCommandCreate(s.State.User.ID, "", resetCmd)
	if err != nil {
		log.Printf("Cannot create '/io-reset' command: %v", err)
	}

	_, err = s.ApplicationCommandCreate(s.State.User.ID, "", showCmd)
	if err != nil {
		log.Printf("Cannot create '/ioshow' command: %v", err)
	}

	log.Println("Commands registered successfully.")
}

// This function is called every time a new interaction is created.
func interactionCreate(s *discordgo.Session, i *discordgo.InteractionCreate) {
	// We only care about application commands (slash commands).
	if i.Type != discordgo.InteractionApplicationCommand {
		return
	}

	commandName := i.ApplicationCommandData().Name

	// Find the voice channel of the user who triggered the command.
	vs, err := findUserVoiceState(DiscordSession, i.GuildID, i.Member.User.ID)
	if err != nil {
		sendEphemeralResponse(s, i, "Error: Could not find you in a voice channel. Please join one to use this command.")
		return
	}

	// Handle the /io command
	if commandName == "io" {
		handleIoCommand(s, i, vs)
	}

	if commandName == "ioshow" {
		handleIoShowCommand(s, i, vs)
	}

	// Handle the /io-reset command
	if commandName == "io-reset" {
		handleIoResetCommand(s, i, vs)
	}
}
