package functions

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/jessicabuzzelli/running-playlist-generator/models"
)

func CreateCombinationPlaylist(c *gin.Context, client *http.Client, songsPlaylistId string, episodesPlaylistId string, outputPlaylistId string, numPodcasts int, numSongsInBetween int) (string, *[]models.PlaylistItemData, error) {

	songCandidates, err := getPlaylistItems(c, client, songsPlaylistId, "song")
	if err != nil {
		return "", nil, err
	}

	episodeCandidates, err := getPlaylistItems(c, client, episodesPlaylistId, "episode")
	if err != nil {
		return "", nil, err
	}

	playlistChoices, err := composePlaylistFromTracks(c, client, songCandidates, episodeCandidates, numPodcasts, numSongsInBetween)
	if err != nil {
		return "", playlistChoices, err
	}

	playlistId, err := createOrResetPlaylist(c, client, outputPlaylistId)
	if err != nil {
		return "", nil, err
	}

	playlistUrl := playlistBaseURL + playlistId

	err = addChoicesToPlaylist(c, client, playlistId, playlistChoices)
	if err != nil {
		return playlistUrl, playlistChoices, err
	}

	return playlistUrl, playlistChoices, nil
}
