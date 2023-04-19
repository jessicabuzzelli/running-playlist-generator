package models

type PlaylistItemData struct {
	URI         string `json:"uri"`
	ID          string `json:"id"`
	Name        string `json:"name"`
	FullyPlayed bool   `json:"fully_played"`
}

func RemovePlaylistItemData(s []PlaylistItemData, i int) []PlaylistItemData {
	length := len(s)
	s[i] = s[length-1]
	return s[:length-1]
}
