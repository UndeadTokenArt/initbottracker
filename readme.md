### How to Run This Project

1.  **Create a Discord Bot:**
    * Go to the [Discord Developer Portal](https://discord.com/developers/applications).
    * Click "New Application". Give it a name (e.g., "InitiativeTracker").
    * Go to the "Bot" tab on the left. Click "Add Bot".
    * Under the bot's username, click "Reset Token" to reveal your bot token. **Keep this secret!**
    * Enable the `SERVER MEMBERS INTENT` and `MESSAGE CONTENT INTENT` under "Privileged Gateway Intents".

2.  **Set Up Your Go Environment:**
    * Make sure you have Go installed on your computer.
    * Create a new folder for your project (e.g., `discord-bot`).
    * Save the Go code above as `main.go` inside that folder.
    * Open a terminal or command prompt in that folder.
    * Initialize Go modules: `go mod init mybot`
    * Get the `discordgo` dependency: `go get github.com/bwmarrin/discordgo`

3.  **Run the Bot:**
    * Set the bot token you got from the developer portal as an environment variable.
        * **On macOS/Linux:** `export DISCORD_BOT_TOKEN="YOUR_TOKEN_HERE"`
        * **On Windows (CMD):** `set DISCORD_BOT_TOKEN="YOUR_TOKEN_HERE"`
        * **On Windows (PowerShell):** `$env:DISCORD_BOT_TOKEN="YOUR_TOKEN_HERE"`
    * Run the bot from your terminal: `go run main.go`
    * You should see messages saying "Bot is now running" and "Web server starting".

4.  **Invite the Bot to Your Server:**
    * In the Developer Portal, go to the "OAuth2" -> "URL Generator" page for your application.
    * Select the `bot` and `applications.commands` scopes.
    * Under "Bot Permissions", select "Send Messages" and "Read Messages/View Channels".
    * Copy the generated URL, paste it into your browser, and invite the bot to your server.

5.  **View the Website:**
    * Save the HTML code above as `index.html`.
    * Open the `index.html` file in your web browser. It will start trying to connect to your bot.

6.  **Use the Command:**
    * Join a voice channel in your Discord server.
    * In any text channel, type `/io` and you should see the command pop up. Use it like `/io 14`.
    * The website should update within a few seconds to show you and anyone else in the channel who has set their initiative. Use `/io-reset` to clear the boa