package functions

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"math/rand"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jessicabuzzelli/running-playlist-generator/models"
	spotify "github.com/zmb3/spotify/v2"
)

const (
	baseURL         = "https://api.spotify.com/v1"
	playlistBaseURL = "https://open.spotify.com/playlist/"

	// hardcoded default inputs
	defaultPlaylistName        = "running playlist -- test"
	defaultPlaylistDescription = ""
	defaultPublic              = false
	defaultCollaborative       = false
	numSongsInBetween          = 4
	numPodcasts                = 3
)

func getPlaylistItems(c *gin.Context, client *http.Client, playlistId string, itemType string) (*[]models.PlaylistItemData, error) {
	nextUrl := fmt.Sprintf("%s/playlists/%s/tracks", baseURL, playlistId)

	var items []spotify.PlaylistItem

	iteration := 0

	for {
		if nextUrl == "" || iteration > 10 {
			break
		}

		resp, err := client.Get(nextUrl)
		if err != nil {
			return nil, fmt.Errorf(fmt.Sprintf("%s request failed!", nextUrl))
		}
		defer resp.Body.Close()

		var playlistPage spotify.PlaylistItemPage
		err = json.NewDecoder(resp.Body).Decode(&playlistPage)
		if err != nil {
			return nil, err
		}

		items = append(items, playlistPage.Items...)

		nextUrl = playlistPage.Next

		iteration += 1
	}

	output := make([]models.PlaylistItemData, len(items))

	for idx, track := range items {
		if (itemType == "both" || itemType == "episode") && track.Track.Episode != nil {
			output[idx] = models.PlaylistItemData{
				URI:         "spotify:episode:" + string(track.Track.Episode.ID),
				ID:          string(track.Track.Episode.ID),
				Name:        track.Track.Episode.Name,
				FullyPlayed: track.Track.Episode.ResumePoint.FullyPlayed,
			}
		}

		if (itemType == "both" || itemType == "song") && track.Track.Track != nil {
			output[idx] = models.PlaylistItemData{
				URI:         "spotify:track:" + string(track.Track.Track.ID),
				ID:          string(track.Track.Track.ID),
				Name:        track.Track.Track.Name,
				FullyPlayed: true,
			}
		}
	}

	return &output, nil
}

func composePlaylistFromTracks(c *gin.Context, client *http.Client, songCandidates *[]models.PlaylistItemData, episodeCandidates *[]models.PlaylistItemData, numPodcasts int, numSongsInBetween int) (*[]models.PlaylistItemData, error) {
	// playlist will follow pattern {song1, ..., songN, podcast1, songN+1, ..., songN+N, podcast2, ...}
	playlist := make([]models.PlaylistItemData, 0, (numPodcasts*numSongsInBetween)+1)

	songs := *songCandidates
	episodes := *episodeCandidates

	if len(songs) < numPodcasts*numSongsInBetween {
		return nil, errors.New(fmt.Sprintf("not enough songs!"))
	}

	if len(episodes) < numPodcasts {
		return nil, errors.New(fmt.Sprintf("not enough podcasts!"))
	}

	for p := 0; p < numPodcasts; p++ {
		for s := 0; s < numSongsInBetween; s++ {
			songIdx := rand.Intn(len(songs))
			song := songs[songIdx]
			playlist = append(playlist, song)

			songs = models.RemovePlaylistItemData(songs, songIdx)
		}

		episodeIdx := rand.Intn(len(episodes))
		episode := episodes[episodeIdx]
		playlist = append(playlist, episode)

		episodes = models.RemovePlaylistItemData(episodes, episodeIdx)
	}

	return &playlist, nil
}

func createOrResetPlaylist(c *gin.Context, client *http.Client, playlistId string) (string, error) {
	if playlistId == "" {
		playlistId, err := createPlaylist(c, client)
		if err != nil {
			return "", err
		} else {
			return playlistId, nil
		}
	} else {
		err := resetPlaylist(c, client, playlistId)
		if err != nil {
			return playlistId, err
		} else {
			return playlistId, nil
		}
	}
}

func createPlaylist(c *gin.Context, client *http.Client) (string, error) {
	userId, err := getUserId(c, client)
	if err != nil {
		return "", err
	}

	url := fmt.Sprintf("%s/users/%s/playlists", baseURL, userId)

	// todo: use inputs from FE
	payload := &models.NewPlaylistPayload{
		Name:          defaultPlaylistName,
		Public:        defaultPublic,
		Collaborative: defaultCollaborative,
		Description:   defaultPlaylistDescription,
	}

	body, _ := json.Marshal(payload)

	resp, err := client.Post(url, "application/json", bytes.NewBuffer(body))
	if err != nil {
		println("failed to post")
		return "", err
	}
	defer resp.Body.Close()

	var playlist spotify.SimplePlaylist
	err = json.NewDecoder(resp.Body).Decode(&playlist)
	if err != nil {
		println("failed to decode new playlist")
		return "", err
	}

	return string(playlist.ID), nil
}

func resetPlaylist(c *gin.Context, client *http.Client, playlistId string) error {
	itemsToRemove, err := getPlaylistItems(c, client, playlistId, "both")
	if err != nil {
		return err
	}

	if len(*itemsToRemove) == 0 {
		println("nothing to remove from playlist " + playlistId)
		return nil
	}

	err = removeItemsFromPlaylist(c, client, playlistId, itemsToRemove)
	if err != nil {
		return err
	}

	return nil
}

func getUserId(c *gin.Context, client *http.Client) (string, error) {
	url := fmt.Sprintf("%s/me", baseURL)

	resp, err := client.Get(url)
	if err != nil {
		fmt.Printf("%s request failed!", url)
	}
	defer resp.Body.Close()

	var user spotify.PrivateUser
	err = json.NewDecoder(resp.Body).Decode(&user)
	if err != nil {
		return "", err
	}

	return string(user.User.ID), nil
}

func removeItemsFromPlaylist(c *gin.Context, client *http.Client, playlistId string, itemsToRemove *[]models.PlaylistItemData) error {
	url := fmt.Sprintf("%s/playlists/%s/tracks", baseURL, playlistId)

	payload := make([]models.URI, 0)

	for _, item := range *itemsToRemove {
		if item.URI != "" {
			payload = append(payload, models.URI{URI: item.URI})
		}
	}

	tracksPayload := make(map[string][]models.URI)
	tracksPayload["tracks"] = payload

	body, _ := json.Marshal(tracksPayload)

	req, err := http.NewRequest("DELETE", url, bytes.NewBuffer(body))
	if err != nil {
		return err
	}

	resp, err := client.Do(req)
	if err != nil {
		return err
	}

	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return errors.New("failed to delete podcast items")
	}

	return nil
}

func addChoicesToPlaylist(c *gin.Context, client *http.Client, playlistId string, playlistChoices *[]models.PlaylistItemData) error {
	url := fmt.Sprintf("%s/playlists/%s/tracks", baseURL, playlistId)

	newItemUris := make([]string, 0)
	for _, choice := range *playlistChoices {
		if choice.URI != "" {
			newItemUris = append(newItemUris, string(choice.URI))
		}
	}

	payload := make(map[string][]string)
	payload["uris"] = newItemUris

	body, err := json.Marshal(payload)
	if err != nil {
		println(err.Error())
		return err
	}

	resp, err := client.Post(url, "application/json", bytes.NewBuffer(body))
	if err != nil {
		println("failed to post")
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 201 {
		println("failed to add items to new playlist")
		return errors.New("failed to add items to playlist")
	}

	return nil
}
