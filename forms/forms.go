package forms

type RunningPlaylistGenerator struct {
	SongPlaylistUrl    string `json:"song_playlist_url" binding:"required"`
	PodcastPlaylistUrl string `json:"podcast_playlist_url" binding:"required"`
	OutputPlaylistUrl  string `json:"output_playlist_url" binding:"required"`
}

type PodcastPlaylistUpdater struct {
	PodcastPlaylistUrl string `json:"podcast_playlist_url" binding:"required"`
}
