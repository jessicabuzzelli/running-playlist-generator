package handlers

import (
	"net/http"
	"os"
	"strconv"

	"github.com/gin-gonic/gin"

	"github.com/jessicabuzzelli/running-playlist-generator/functions"
)

func GenerateRunningPlaylist(c *gin.Context) {

	// todo : pull from ctx
	outputPlaylistId := os.Getenv("outputPlaylistId")
	episodesPlaylistId := os.Getenv("episodesPlaylistId")
	songsPlaylistId := os.Getenv("songsPlaylistId")
	numPodcasts, _ := strconv.Atoi(os.Getenv("numPodcasts"))
	numSongsInBetween, _ := strconv.Atoi(os.Getenv("numSongsInBetween"))

	playlistUrl, _, err := functions.CreateCombinationPlaylist(c, client, songsPlaylistId, episodesPlaylistId, outputPlaylistId, numPodcasts, numSongsInBetween)
	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}

	if err != nil {
		c.String(http.StatusInternalServerError, "received playlist track choices in an invalid format")
		return
	} else {
		c.String(http.StatusOK, playlistUrl)
		return
	}
}

func UpdatePodcastPlaylist(c *gin.Context) {

	// todo : pull from ctx
	playlistId := os.Getenv("episodesPlaylistId")

	url, err := functions.RemovePlayedItemsFromPlaylist(c, client, playlistId)
	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	} else {
		c.String(http.StatusOK, url)
		return
	}
}
