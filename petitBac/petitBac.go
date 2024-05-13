package petitBac

import (
	"fmt"
	"html/template"
	"log"
	"math/rand"
	"net/http"
	"time"
)

type GameResponse struct {
	Letter     rune
	Artist     string
	Album      string
	MusicGroup string
	Instrument string
	Featuring  string
}

var responses []GameResponse
var htmlContent []byte
var letters []rune
var originalLetters = []rune("ABCDEFGHIJKLMNOPQRSTUVWXYZ")
var templates *template.Template

func init() {
	var err error
	// Charger les templates directement par leur nom si elles sont dans le même répertoire
	templates, err = template.ParseFiles("petit-bac.html", "petit-bac-answers.html")
	if err != nil {
		log.Fatal("Error loading templates:", err)
	}
	resetLetters()
}

// Fonction pour réinitialiser la liste des lettres
func resetLetters() {
	letters = make([]rune, len(originalLetters))
	copy(letters, originalLetters)
}

// Fonction pour obtenir une lettre aléatoire
func getRandomLetter() rune {
	if len(letters) == 0 {
		resetLetters()
	}

	rand.Seed(time.Now().UnixNano())
	index := rand.Intn(len(letters))
	letter := letters[index]
	letters = append(letters[:index], letters[index+1:]...)
	return letter
}

// Supposons que vous ayez une méthode pour récupérer la lettre actuellement utilisée
func getCurrentLetter() rune {
	if len(letters) > 0 {
		return letters[len(letters)-1]
	}
	return ' '
}

func PetitBacHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		newResponse := GameResponse{
			Letter:     getCurrentLetter(),
			Artist:     r.FormValue("artist"),
			Album:      r.FormValue("album"),
			MusicGroup: r.FormValue("musicGroup"),
			Instrument: r.FormValue("instrument"),
			Featuring:  r.FormValue("featuring"),
		}
		responses = append(responses, newResponse)
		http.Redirect(w, r, "/petit-bac-answers", http.StatusSeeOther)
		return
	}

	currentLetter := getRandomLetter()
	if err := templates.ExecuteTemplate(w, "petit-bac.html", map[string]interface{}{"Letter": string(currentLetter)}); err != nil {
		http.Error(w, "Error rendering template: "+err.Error(), http.StatusInternalServerError)
	}
}

func AnswersHandler(w http.ResponseWriter, r *http.Request) {
	if len(responses) > 0 {
		lastResponse := responses[len(responses)-1]
		if err := templates.ExecuteTemplate(w, "petit-bac-answers.html", lastResponse); err != nil {
			http.Error(w, "Error rendering template: "+err.Error(), http.StatusInternalServerError)
		}
	} else {
		http.Error(w, "No answers to display", http.StatusNotFound)
	}
}

func ValidateAnswersHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Redirect(w, r, "/petit-bac", http.StatusSeeOther)
		return
	}

	score := 0
	if r.FormValue("artistCorrect") == "true" {
		score++
	}
	if r.FormValue("albumCorrect") == "true" {
		score++
	}
	if r.FormValue("musicGroupCorrect") == "true" {
		score++
	}
	if r.FormValue("instrumentCorrect") == "true" {
		score++
	}
	if r.FormValue("featuringCorrect") == "true" {
		score++
	}

	result := "Sorry, you did not score enough."
	if score >= 3 {
		result = "Congratulations, you win a point!"
	}

	// Modifier ici pour inclure un lien direct vers la page du jeu
	fmt.Fprintf(w, `<html><body><p>%s</p><a href='/petit-bac'>Try again</a></body></html>`, result)
}
