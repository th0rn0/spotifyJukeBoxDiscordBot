package main

import "github.com/zmb3/spotify/v2"

type SearchOutput struct {
	Tracks []Track
}

type Track struct {
	Name   string `json:"name"`
	Artist string `json:"artist"`
	ID     string `json:"id"`
}

// API Requests
type APIAddTrack struct {
	URI spotify.URI `json:"uri"`
}

type APIVolume struct {
	Volume int `json:"volume"`
}

// API Responses
type APIResponseSearchOutput struct {
	Playlists []APIResponsePlaylistSearchOutput `json:"playlist"`
	Artists   []APIResponseArtistSearchOutput   `json:"artist"`
	Tracks    []APIResponseTrackSearchOutput    `json:"track"`
}

type APIResponsePlaylistSearchOutput struct {
	Name string `json:"name"`
	ID   string `json:"id"`
}

type APIResponseArtistSearchOutput struct {
	Name string `json:"name"`
	ID   string `json:"id"`
}

type APIResponseTrackSearchOutput struct {
	Name   string          `json:"name"`
	Artist string          `json:"artist"`
	ID     string          `json:"id"`
	Images []spotify.Image `json:"images"`
}
