package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"net/url"
	"strings"
	"time"
)

const clientID = "dcb92b47b4fd450094f6e91c12bd1e4d"
const clientSecret = "fa6ef02fc29e4def8b8bf60bdff4ea75"
const tokenURL = "https://accounts.spotify.com/api/token"
const playlistID = "1vCUdlD8Ic1KyEMctetRbU"
const SpotifyAPIBase = "https://api.spotify.com/v1"

func getAccessToken() (string, error) {
	client := &http.Client{}
	data := url.Values{}
	data.Set("grant_type", "client_credentials")

	req, err := http.NewRequest("POST", tokenURL, strings.NewReader(data.Encode()))
	if err != nil {
		return "", err
	}

	req.SetBasicAuth(clientID, clientSecret)
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	var result map[string]interface{}
	if err := json.Unmarshal(body, &result); err != nil {
		return "", err
	}

	token, ok := result["access_token"].(string)
	if !ok {
		return "", fmt.Errorf("access token not found in the response")
	}

	return token, nil
}

func getRandomTrackFromPlaylist(accessToken, playlistID string) (string, string, error) {
	client := &http.Client{}
	req, err := http.NewRequest("GET", SpotifyAPIBase+"/playlists/"+playlistID+"/tracks", nil)
	if err != nil {
		return "", "", err
	}

	req.Header.Add("Authorization", "Bearer "+accessToken)
	resp, err := client.Do(req)
	if err != nil {
		return "", "", err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", "", err
	}

	var result map[string]interface{}
	if err := json.Unmarshal(body, &result); err != nil {
		return "", "", fmt.Errorf("error parsing JSON response: %s, response body: '%s'", err, string(body))
	}

	items, ok := result["items"].([]interface{})
	if !ok || len(items) == 0 {
		return "", "", fmt.Errorf("no tracks found in playlist")
	}

	rand.Seed(time.Now().UnixNano())
	randomIndex := rand.Intn(len(items))
	track := items[randomIndex].(map[string]interface{})["track"].(map[string]interface{})

	trackName, ok := track["name"].(string)
	if !ok {
		return "", "", fmt.Errorf("track name not found")
	}

	artists := track["artists"].([]interface{})
	if len(artists) == 0 {
		return "", "", fmt.Errorf("no artists found for track")
	}

	artistName, ok := artists[0].(map[string]interface{})["name"].(string)
	if !ok {
		return "", "", fmt.Errorf("artist name not found")
	}

	return trackName, artistName, nil
}

func getLyricsFromMusixmatch(trackName, artistName, apiKey string) (string, error) {
	baseURL := "https://api.musixmatch.com/ws/1.1/"
	endpoint := "matcher.lyrics.get"

	// Build the query URL with the necessary parameters
	url := fmt.Sprintf("%s%s?format=json&apikey=%s&q_track=%s&q_artist=%s", baseURL, endpoint, apiKey, url.QueryEscape(trackName), url.QueryEscape(artistName))

	resp, err := http.Get(url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("failed to get lyrics, status code: %d", resp.StatusCode)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	var data map[string]interface{}
	if err := json.Unmarshal(body, &data); err != nil {
		return "", err
	}

	lyrics := data["message"].(map[string]interface{})["body"].(map[string]interface{})["lyrics"].(map[string]interface{})["lyrics_body"].(string)

	return lyrics, nil
}

func main() {
	accessToken, err := getAccessToken()
	if err != nil {
		log.Fatalf("Failed to get access token: %s", err)
	}

	trackName, artistName, err := getRandomTrackFromPlaylist(accessToken, playlistID)
	if err != nil {
		log.Fatalf("Error: %s", err)
	}

	// Get song lyrics from the Musixmatch API
	lyrics, err := getLyricsFromMusixmatch(trackName, artistName, "fcc277ce6c9bd4d25476e2107fffec18")
	if err != nil {
		log.Printf("Failed to get lyrics: %s", err)
	} else {
		fmt.Printf("Lyrics for %s by %s:\n%s\n", trackName, artistName, lyrics)
	}
}
