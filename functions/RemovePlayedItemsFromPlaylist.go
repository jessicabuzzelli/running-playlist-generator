package functions

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jessicabuzzelli/running-playlist-generator/models"
)

func RemovePlayedItemsFromPlaylist(c *gin.Context, client *http.Client, playlistId string) (string, error) {
	url := playlistBaseURL + playlistId

	episodes, err := getPlaylistItems(c, client, playlistId, "both")
	if err != nil {
		return "", err
	}

	itemsToRemove := make([]models.PlaylistItemData, 0)
	for _, episode := range *episodes {
		if episode.FullyPlayed {
			itemsToRemove = append(itemsToRemove, episode)
		}
	}

	if len(itemsToRemove) == 0 {
		fmt.Print("nothing to remove from playlist " + playlistId)
		return url, nil
	}

	err = removeItemsFromPlaylist(c, client, playlistId, &itemsToRemove)
	if err != nil {
		return "", err
	}

	return url, nil
}
