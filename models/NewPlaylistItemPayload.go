package models

type NewPlaylistItemPayload struct {
	URIs     []string `json:"uris"`
	Position int      `json:"position"`
}
