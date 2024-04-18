package blindtest

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
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

var httpClient = &http.Client{}

func getAccessToken() (string, error) {
	data := url.Values{}
	data.Set("grant_type", "client_credentials")

	req, err := http.NewRequest("POST", tokenURL, strings.NewReader(data.Encode()))
	if err != nil {
		return "", fmt.Errorf("failed to create request: %v", err)
	}

	req.SetBasicAuth(clientID, clientSecret)
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	resp, err := httpClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to send request: %v", err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read response body: %v", err)
	}

	var result map[string]interface{}
	if err := json.Unmarshal(body, &result); err != nil {
		return "", fmt.Errorf("JSON unmarshal error: %v, body: %s", err, body)
	}

	token, ok := result["access_token"].(string)
	if !ok {
		return "", fmt.Errorf("access token not found in the response")
	}

	return token, nil
}

func getRandomTrackFromPlaylist(accessToken, playlistID string) (string, string, string, error) {
	client := &http.Client{}
	req, err := http.NewRequest("GET", SpotifyAPIBase+"/playlists/"+playlistID+"/tracks", nil)
	if err != nil {
		return "", "", "", err
	}

	req.Header.Add("Authorization", "Bearer "+accessToken)
	resp, err := client.Do(req)
	if err != nil {
		return "", "", "", err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", "", "", err
	}

	var result map[string]interface{}
	if err := json.Unmarshal(body, &result); err != nil {
		return "", "", "", fmt.Errorf("error parsing JSON response: %s, response body: '%s'", err, string(body))
	}

	items, ok := result["items"].([]interface{})
	if !ok || len(items) == 0 {
		return "", "", "", fmt.Errorf("no tracks found in playlist")
	}

	rand.Seed(time.Now().UnixNano())
	randomIndex := rand.Intn(len(items))
	track := items[randomIndex].(map[string]interface{})["track"].(map[string]interface{})

	trackName, ok := track["name"].(string)
	if !ok {
		return "", "", "", fmt.Errorf("track name not found")
	}

	artists := track["artists"].([]interface{})
	if len(artists) == 0 {
		return "", "", "", fmt.Errorf("no artists found for track")
	}

	artistName, ok := artists[0].(map[string]interface{})["name"].(string)
	if !ok {
		return "", "", "", fmt.Errorf("artist name not found")
	}

	trackURI, ok := track["uri"].(string)
	if !ok {
		return "", "", "", fmt.Errorf("track URI not found")
	}

	return trackName, artistName, trackURI, nil
}

func BlindTestHandler(w http.ResponseWriter, r *http.Request) {
	// Obtention du token d'accès
	accessToken, err := getAccessToken()
	if err != nil {
		http.Error(w, "Failed to get access token", http.StatusInternalServerError)
		return
	}

	// Obtention d'une piste aléatoire de la playlist
	trackName, artistName, trackURI, err := getRandomTrackFromPlaylist(accessToken, playlistID)
	if err != nil {
		http.Error(w, "Failed to get random track", http.StatusInternalServerError)
		fmt.Println(trackName, artistName)
		return
	}

	// Extraire l'ID de la piste depuis l'URI Spotify
	splitURI := strings.Split(trackURI, ":")
	if len(splitURI) < 3 {
		http.Error(w, "Invalid Spotify URI", http.StatusInternalServerError)
		return
	}
	trackID := splitURI[2]

	// Génération du contenu HTML
	htmlContent := fmt.Sprintf(`<!DOCTYPE html>
<html lang="fr">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Music player</title>
    <link rel="stylesheet" href="./static/styles.css">
</head>
<body>
    <iframe style="border-radius:12px" src="https://open.spotify.com/embed/track/%s?utm_source=generator" width="100%%" height="352" frameBorder="0" allowfullscreen="" allow="autoplay; clipboard-write; encrypted-media; fullscreen; picture-in-picture" loading="lazy"></iframe>
</body>
</html>`, trackID)

	w.Header().Set("Content-Type", "text/html")
	w.Write([]byte(htmlContent))
}
