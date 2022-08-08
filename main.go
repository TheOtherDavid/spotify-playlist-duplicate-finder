package main

import (
	"encoding/csv"
	"fmt"
	"os"
	"strings"
)

func main() {
	//Get Spotify auth
	fmt.Println("Refreshing access token.")
	accessToken, err := RefreshSpotifyAuth()
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("Access token refreshed.")
	fmt.Println(accessToken)

	//Get Spotify playlist by ID
	userId := os.Getenv("USER_ID")
	playlistId := os.Args[1]
	if playlistId == "" {
		fmt.Println("Playlist ID not given as an argument. Please specify a Spotify Playlist ID.")
	}

	playlistTracks, err := GetTracksForPlaylist(playlistId, userId, accessToken)
	if err != nil {
		fmt.Println(err)
		return
	}

	potentialDuplicates, err := GetPotentialDuplicateTracks(playlistTracks)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(potentialDuplicates)

	filename := os.Args[2]
	if filename == "" {
		filename = "duplicate-spotify-track"
	}

	//Take that list of tracks and write the band and track name to a file
	WriteDuplicateTracksToFile(potentialDuplicates, filename)
	//Use the file to manually take care of the duplicates, to avoid any automation disasters.
}

func GetPotentialDuplicateTracks(playlistTracks []SpotifyTrack) (potentialDuplicateTracks []SpotifyTrack, err error) {
	artistTrackMap, err := GetArtistTrackMap(playlistTracks)
	if err != nil {
		return nil, err
	}

	potentialDuplicateTracks, err = GetPotentialDuplicateTracksForArtistTrackMap(artistTrackMap)
	if err != nil {
		return nil, err
	}

	return potentialDuplicateTracks, nil
}

func GetArtistTrackMap(playlistTracks []SpotifyTrack) (artistTrackMap map[string][]SpotifyTrack, err error) {
	//Duplicates have to share the same artist.
	//Separate the entire list of songs into lists of each artist's songs, to cut down on the number of comparisons required
	artistTrackMap = make(map[string][]SpotifyTrack)
	for _, track := range playlistTracks {
		//If track has multiple artists, add them both? I think this is the best. False positives are better than false negatives.
		for _, artist := range track.Artists {
			artistName := artist.Name
			//If artist is already in map, add the track to the track list
			if val, ok := artistTrackMap[artistName]; ok {
				artistTrackMap[artistName] = append(val, track)
			} else {
				//If artist is not in map, add artist to map with track
				tracks := []SpotifyTrack{track}
				artistTrackMap[artistName] = tracks
			}
		}
	}
	return artistTrackMap, nil
}

func GetPotentialDuplicateTracksForArtistTrackMap(artistTrackMap map[string][]SpotifyTrack) (potentialDuplicateTracks []SpotifyTrack, err error) {
	//Now we have to do the substring search.
	//For each artist, look at the list of tracks
	potentialDuplicateTracks = []SpotifyTrack{}
	for _, tracks := range artistTrackMap {

		for i, trackA := range tracks {
			//For each track, compare it to the rest of the tracks in the list
			trackAName := trackA.Name
			for _, trackB := range tracks[i+1:] {
				//NOW we finally compare the substrings. If trackA's name appears within trackB's name, add them BOTH to the potential dupes.
				trackBName := trackB.Name
				if strings.Contains(trackBName, trackAName) {
					potentialDuplicateTracks = append(potentialDuplicateTracks, trackA)
					potentialDuplicateTracks = append(potentialDuplicateTracks, trackB)
				}
			}
		}
	}
	return potentialDuplicateTracks, nil
}

func WriteDuplicateTracksToFile(potentialDuplicateTracks []SpotifyTrack, filename string) (err error) {
	file, err := os.Create(filename + ".csv")
	if err != nil {
		return err
	}
	defer file.Close()
	w := csv.NewWriter(file)
	defer w.Flush()

	// Using WriteAll
	var data [][]string
	for _, track := range potentialDuplicateTracks {
		row := []string{track.Artists[0].Name, track.Album.Name, track.Name}
		data = append(data, row)
	}
	w.WriteAll(data)
	return nil
}
