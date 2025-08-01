package webserver

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/undeadtokenart/initbottracker/botcommands"
)

func StartWebServer() {
	router := gin.Default()

	// serve index.html from the templates directory
	router.LoadHTMLGlob("templates/*.html")

	// Define the route for the index page
	router.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.html", gin.H{})
	})

	// Define the route for the initiative endpoint using Gin
	router.GET("/initiative", botcommands.InitiativeHandler)

	log.Println("Web server starting on http://localhost:8081")
	if err := router.Run(":8081"); err != nil {
		log.Fatalf("Failed to start web server: %v", err)
	}
}
