package blindtest

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/rand"
	"net/http"
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
	data := "grant_type=client_credentials"
	req, err := http.NewRequest("POST", tokenURL, strings.NewReader(data))
	if err != nil {
		return "", fmt.Errorf("failed to create request: %v", err)
	}
	req.SetBasicAuth(clientID, clientSecret)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

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

func getRandomTrackFromPlaylist(accessToken string) (string, string, string, error) {
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/playlists/%s/tracks", SpotifyAPIBase, playlistID), nil)
	if err != nil {
		return "", "", "", err
	}
	req.Header.Set("Authorization", "Bearer "+accessToken)

	resp, err := httpClient.Do(req)
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
		return "", "", "", fmt.Errorf("error parsing JSON response: %v", err)
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
	// Vérifier si la méthode HTTP est POST (traitement du formulaire)
	if r.Method == http.MethodPost {
		// Récupérer les réponses soumises par l'utilisateur depuis le formulaire
		submittedTrackName := r.FormValue("songName")
		submittedArtistName := r.FormValue("artistName")
		correctTrackName := r.FormValue("trackName")
		correctArtistName := r.FormValue("artistNameCorrect")

		// Vérifier si les réponses soumises sont correctes
		isCorrect := submittedTrackName == correctTrackName && submittedArtistName == correctArtistName

		// Construire le message de résultat
		var resultMessage string
		if isCorrect {
			resultMessage = "Bonne réponse!"
		} else {
			resultMessage = fmt.Sprintf("Mauvaise réponse. La chanson était: %s par %s", correctTrackName, correctArtistName)
		}

		// Afficher le résultat dans la même page
		htmlContent := fmt.Sprintf(`<!DOCTYPE html>
        <html lang="fr">
        <head>
            <meta charset="UTF-8">
            <meta name="viewport" content="width=device-width, initial-scale=1.0">
            <title>Résultat Blind Test Spotify</title>
            <link rel="stylesheet" href="/static/blind-test.css">
        </head>
        <body>
            <h1>Résultat du Blind Test</h1>
            <p>%s</p>
            <a href="/blind-test">Retour au jeu</a>
        </body>
        </html>`, resultMessage)

		// Envoyer la réponse HTML avec le résultat
		w.Header().Set("Content-Type", "text/html")
		w.Write([]byte(htmlContent))

		return
	}

	// Si ce n'est pas une requête POST, afficher le formulaire normal
	accessToken, err := getAccessToken()
	if err != nil {
		http.Error(w, "Failed to get access token", http.StatusInternalServerError)
		return
	}

	trackName, artistName, trackURI, err := getRandomTrackFromPlaylist(accessToken)
	if err != nil {
		http.Error(w, "Failed to get random track", http.StatusInternalServerError)
		return
	}

	splitURI := strings.Split(trackURI, ":")
	if len(splitURI) < 3 {
		http.Error(w, "Invalid Spotify URI", http.StatusInternalServerError)
		return
	}
	trackID := splitURI[2]

	htmlContent := fmt.Sprintf(`<!DOCTYPE html>
    <html lang="fr">
    <head>
        <meta charset="UTF-8">
        <meta name="viewport" content="width=device-width, initial-scale=1.0">
        <title>Blind Test Spotify</title>
        <link rel="stylesheet" href="/static/blind-test.css">
    </head>
    <body>
        <h1>Blind Test Spotify</h1>
        <iframe id="spotifyPlayer" src="https://open.spotify.com/embed/track/%s" width="300" height="80" frameborder="0" allowtransparency="true" allow="encrypted-media"></iframe>
        <form method="post">
            <label for="songName">Nom de la chanson :</label>
            <input type="text" id="songName" name="songName" required><br>
            <label for="artistName">Nom de l'artiste :</label>
            <input type="text" id="artistName" name="artistName" required><br>
            <input type="hidden" name="trackName" value="%s">
            <input type="hidden" name="artistNameCorrect" value="%s">
            <input type="submit" value="Valider">
        </form>
    </body>
    </html>`, trackID, trackName, artistName)

	w.Header().Set("Content-Type", "text/html")
	w.Write([]byte(htmlContent))
}

func CheckAnswerHandler(w http.ResponseWriter, r *http.Request) {
	// Vérifier si la méthode HTTP est POST
	if r.Method != http.MethodPost {
		http.Error(w, "Méthode non autorisée", http.StatusMethodNotAllowed)
		return
	}

	// Récupérer la réponse de l'utilisateur depuis le formulaire
	userAnswer := r.FormValue("userAnswer")

	// Liste des réponses correctes possibles
	correctAnswers := []string{"Je danse le Mia", "Lettre ", "Fenêtre Sur Rue", "Dans ma rue", "High for the Chronic", "Samuraï", "Samurai", "Ma Benz", "ma benz", "Petit frère", "Qui est l'exemple", "shurikn", "shurik'n", "Shurik'n", "Supreme NTM", "supreme NTM", "NTM", "hugo tsr", "Hugo TSR"}

	// Convertir la réponse de l'utilisateur et les réponses correctes en minuscules pour la comparaison
	userAnswerLower := strings.ToLower(userAnswer)

	// Vérifier si la réponse de l'utilisateur correspond à l'une des réponses correctes
	var correct bool
	for _, correctAnswer := range correctAnswers {
		if userAnswerLower == strings.ToLower(correctAnswer) {
			correct = true
			break
		}
	}

	// Envoyer une réponse en fonction de la vérification
	if correct {
		// Envoyer une réponse de succès si la réponse est correcte
		fmt.Fprintln(w, "Bravo, vous avez deviné la bonne chanson !")
	} else {
		// Envoyer une réponse d'échec si la réponse est incorrecte
		fmt.Fprintln(w, "Désolé, votre réponse est incorrecte.")
	}
}
