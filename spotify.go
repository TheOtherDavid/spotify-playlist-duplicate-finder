package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strconv"
)

type SpotifyRefreshTokenResponse struct {
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
	ExpiresIn   int    `json:"expires_in"`
	Scope       string `json:"scope"`
}

func RefreshSpotifyAuth() (string, error) {
	clientId := os.Getenv("SPOTIFY_ID")
	clientSecret := os.Getenv("SPOTIFY_SECRET")
	refreshToken := os.Getenv("REFRESH_TOKEN")
	grantType := "refresh_token"
	url := "https://accounts.spotify.com/api/token?client_id=" +
		clientId + "&client_secret=" + clientSecret + "&refresh_token=" +
		refreshToken + "&grant_type=" + grantType

	response, err := http.Post(url, "application/x-www-form-urlencoded", nil)
	if err != nil {
		fmt.Println("Oh no, error.")
	}
	defer response.Body.Close()

	var responseBody SpotifyRefreshTokenResponse
	err = json.NewDecoder(response.Body).Decode(&responseBody)
	if err != nil {
		fmt.Println("Oh no, error.")
		return "", err
	}

	accessToken := responseBody.AccessToken
	return accessToken, nil
}

type SpotifyPlaylistResponse struct {
	Items  []SpotifyPlaylistTrackItem `json:"items"`
	Limit  int                        `json:"limit"`
	Offset int                        `json:"offset"`
	Total  int                        `json:"total"`
}

type SpotifyPlaylistTrackItem struct {
	Track SpotifyTrack `json:"track"`
}
type SpotifyTrack struct {
	Album       SpotifyAlbum        `json:"album"`
	Artists     []SpotifyArtistItem `json:"artists"`
	Id          string              `json:"id"`
	Name        string              `json:"name"`
	Popularity  int                 `json:"popularity"`
	TrackNumber int                 `json:"track_number"`
}

type SpotifyAlbum struct {
	Id          string `json:"id"`
	Name        string `json:"name"`
	ReleaseDate string `json:"release_date"`
}

type SpotifyArtistItem struct {
	Genres     []string `json:"genres"`
	Id         string   `json:"id"`
	Name       string   `json:"name"`
	Popularity int      `json:"popularity"`
	Type       string   `json:"type"`
}

func GetTracksForPlaylist(playlistId string, userId string, accessToken string) (tracks []SpotifyTrack, err error) {
	//Call the whole thing in a loop until the offset is longer than the remaining total of songs
	offset := 0
	total := 0
	limit := 100
	for {
		//Call get Playlist for Playlist ID
		url := "https://api.spotify.com/v1/users/" + userId + "/playlists/" + playlistId + "/tracks?offset=" + strconv.Itoa(offset) + "&limit=" + strconv.Itoa(limit)

		client := http.Client{}
		req, err := http.NewRequest("GET", url, nil)
		if err != nil {
			fmt.Println("Oh no, error.")
			return nil, err
		}

		req.Header = http.Header{
			"Authorization": {"Bearer " + accessToken},
		}

		response, err := client.Do(req)
		if err != nil {
			fmt.Println("Oh no, error.")
			return nil, err
		}
		defer response.Body.Close()

		var responseBody SpotifyPlaylistResponse
		err = json.NewDecoder(response.Body).Decode(&responseBody)
		if err != nil {
			fmt.Println("Oh no, error.")
			return nil, err
		}

		offset = responseBody.Offset
		limit = responseBody.Limit
		total = responseBody.Total
		//Get tracks from response body
		for _, trackItems := range responseBody.Items {
			tracks = append(tracks, trackItems.Track)
		}

		//If the remaining records are less than the limit, make the off
		if (total - offset) < limit {
			break
		}
		offset = offset + limit

		if offset > total {
			break
		}
	}

	return tracks, nil
}
