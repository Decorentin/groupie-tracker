package blindtest

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/rand"
	"net/http"
	"net/url"
	"path/filepath"
	"strings"
	"text/template"
	"time"

	"github.com/gosimple/slug"
)

const clientID = "dcb92b47b4fd450094f6e91c12bd1e4d"
const clientSecret = "fa6ef02fc29e4def8b8bf60bdff4ea75"
const tokenURL = "https://accounts.spotify.com/api/token"
const playlistID = "1vCUdlD8Ic1KyEMctetRbU"
const SpotifyAPIBase = "https://api.spotify.com/v1"

var selectedSongTitle string

var httpClient = &http.Client{}

type BlindTestStruct struct {
	TrackID         string
	TrackName       string
	ArtistName      string
	Token           string
	TrackPreviewURL string // Ajout de l'URL de prévisualisation de la piste
}

var bt *BlindTestStruct

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

func getRandomTrackFromPlaylist(accessToken, playlistID string) (string, string, string, string, error) {
	client := &http.Client{}
	req, err := http.NewRequest("GET", SpotifyAPIBase+"/playlists/"+playlistID+"/tracks", nil)
	if err != nil {
		return "", "", "", "", err
	}

	req.Header.Add("Authorization", "Bearer "+accessToken)
	resp, err := client.Do(req)
	if err != nil {
		return "", "", "", "", err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", "", "", "", err
	}

	var result map[string]interface{}
	if err := json.Unmarshal(body, &result); err != nil {
		return "", "", "", "", fmt.Errorf("error parsing JSON response: %s, response body: '%s'", err, string(body))
	}

	items, ok := result["items"].([]interface{})
	if !ok || len(items) == 0 {
		return "", "", "", "", fmt.Errorf("no tracks found in playlist")
	}

	rand.Seed(time.Now().UnixNano())
	randomIndex := rand.Intn(len(items))
	track := items[randomIndex].(map[string]interface{})["track"].(map[string]interface{})

	trackName, ok := track["name"].(string)
	if !ok {
		return "", "", "", "", fmt.Errorf("track name not found")
	}

	artists := track["artists"].([]interface{})
	if len(artists) == 0 {
		return "", "", "", "", fmt.Errorf("no artists found for track")
	}

	artistName, ok := artists[0].(map[string]interface{})["name"].(string)
	if !ok {
		return "", "", "", "", fmt.Errorf("artist name not found")
	}

	trackURI, ok := track["uri"].(string)
	if !ok {
		return "", "", "", "", fmt.Errorf("track URI not found")
	}

	trackPreviewURL, ok := track["preview_url"].(string)
	if !ok {
		return "", "", "", "", fmt.Errorf("track preview URL not found")
	}

	return trackName, artistName, trackURI, trackPreviewURL, nil
}

func CheckAnswerHandler(w http.ResponseWriter, r *http.Request) {
	// Vérifier si la méthode HTTP est POST
	if r.Method != http.MethodPost {
		http.Error(w, "Méthode non autorisée", http.StatusMethodNotAllowed)
		return
	}

	// Récupérer la réponse de l'utilisateur depuis le formulaire
	userAnswer := r.FormValue("userAnswer")

	// Vérifier si la réponse de l'utilisateur est vide
	if userAnswer == "" {
		// Envoyer un message d'erreur si la réponse est vide
		http.Error(w, "La réponse de l'utilisateur ne peut pas être vide", http.StatusBadRequest)
		return
	}

	// Convertir la réponse de l'utilisateur et le titre de la chanson sélectionnée aléatoirement en minuscules sans accents
	userAnswerLower := removeAccents(strings.ToLower(userAnswer))
	selectedSongTitleLower := removeAccents(strings.ToLower(selectedSongTitle))

	// Vérifier si la réponse de l'utilisateur correspond au titre de la chanson sélectionné aléatoirement
	if userAnswerLower == selectedSongTitleLower {
		// Envoyer une réponse de succès si la réponse est correcte
		fmt.Fprintln(w, "Bravo, vous avez deviné la bonne chanson !")
	} else {
		// Envoyer une réponse d'échec si la réponse est incorrecte
		fmt.Fprintln(w, "Désolé, votre réponse est incorrecte.")
	}
}

func removeAccents(s string) string {
	return slug.Make(s)
}

func BlindTestHandler(w http.ResponseWriter, r *http.Request) {
	// Obtention du token et de la piste comme précédemment
	accessToken, err := getAccessToken()
	if err != nil {
		http.Error(w, "Failed to get access token", http.StatusInternalServerError)
		return
	}
	trackName, artistName, trackURI, trackPreviewURL, err := getRandomTrackFromPlaylist(accessToken, playlistID)
	if err != nil {
		http.Error(w, "Failed to get random track", http.StatusInternalServerError)
		return
	}
	splitURI := strings.Split(trackURI, ":")
	trackID := splitURI[2]

	// Stockez le titre de la chanson sélectionnée aléatoirement dans la variable globale
	selectedSongTitle = trackName

	bt = &BlindTestStruct{
		TrackID:         trackID,
		TrackName:       trackName,
		ArtistName:      artistName,
		Token:           accessToken,
		TrackPreviewURL: trackPreviewURL, // Ajout de l'URL de prévisualisation de la piste
	}

	// Charger le template
	tmplPath := filepath.Join("blind-test.html")
	tmpl, err := template.ParseFiles(tmplPath)
	if err != nil {
		http.Error(w, "Error loading template", http.StatusInternalServerError)
		return
	}

	// Générer la réponse HTML avec le template
	w.Header().Set("Content-Type", "text/html")
	tmpl.Execute(w, bt)
}
