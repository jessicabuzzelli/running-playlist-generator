package main

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/thinkerou/favicon"

	"github.com/jessicabuzzelli/running-playlist-generator/handlers"
)

var (
	ctx = context.Background()
)

func loadEnv() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatalf(fmt.Sprintf("failed to load .env file: %s", err.Error()))
	}
}

func main() {
	loadEnv()

	router := gin.Default()
	router.SetTrustedProxies(nil)

	router.LoadHTMLGlob("templates/*")

	router.Use(favicon.New("assets/favicon.ico"))

	public := router.Group("/")

	public.GET("/", func(c *gin.Context) {
		c.HTML(
			http.StatusOK,
			"index.tmpl",
			gin.H{},
		)
	})

	public.GET("/login", handlers.Login)

	public.GET("/callback", handlers.OAuthCallback)

	public.GET("/loginFailed", func(c *gin.Context) {
		c.HTML(
			http.StatusOK,
			"failure.tmpl",
			gin.H{},
		)
	})

	private := router.Group("/app")

	private.GET("/loginSuccessful", func(c *gin.Context) {
		c.Redirect(http.StatusSeeOther, "/app/home")
	})

	private.GET("/home", func(c *gin.Context) {
		c.HTML(
			http.StatusOK,
			"home.tmpl",
			gin.H{},
		)
	})

	private.GET("/generateRunningPlaylist", handlers.GenerateRunningPlaylist)

	private.GET("/updatePodcastPlaylist", handlers.UpdatePodcastPlaylist)

	router.Run(":8080")

	fmt.Print("http://localhost:8080/")
}
