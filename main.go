package main

import (
	"fmt"
	"log"

	"os"
	"os/signal"

	"syscall"

	"github.com/undeadtokenart/initbottracker/botcommands"
	"github.com/undeadtokenart/initbottracker/webserver"
)

func main() {
	// Create a new Discord session.
	botcommands.CreateDiscordSession()

	// Start the web server in a separate goroutine.
	go webserver.StartWebServer()
	// Log that the bot is ready.
	log.Println("Bot is ready and web server is running.")

	// Wait here until CTRL-C or other term signal is received.
	fmt.Println("Bot is now running. Press CTRL-C to exit.")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-sc

	// Cleanly close down the Discord session.
	fmt.Println("Bot shutting down.")
}
