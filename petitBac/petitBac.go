package petitBac

import (
	"fmt"
	"html/template"
	"log"
	"math/rand"
	"net/http"
	"strconv"
	"time"
)

// GameResponse stores the data for each game response, including the current letter and user inputs.
type GameResponse struct {
	Letter     rune
	Artist     string
	Album      string
	MusicGroup string
	Instrument string
	Featuring  string
}

var responses []GameResponse                               // Stores all game responses for the session
var letters []rune                                         // Dynamic slice to keep track of available letters
var originalLetters = []rune("ABCDEFGHIJKLMNOPQRSTUVWXYZ") // Original set of all letters
var templates *template.Template                           // Holds the loaded HTML templates

// init is called when the package is initialized. It sets up the templates and resets the letters.
func init() {
	var err error
	// Load HTML templates from files
	templates, err = template.ParseFiles("petit-bac.html", "petit-bac-answers.html")
	if err != nil {
		log.Fatal("Error loading templates:", err)
	}
	resetLetters() // Initialize the letters slice
}

// resetLetters reinitializes the letters slice to its original full alphabet state.
func resetLetters() {
	letters = make([]rune, len(originalLetters))
	copy(letters, originalLetters)
}

// getRandomLetter selects a random letter from the available set, removes it, and returns it.
func getRandomLetter() rune {
	if len(letters) == 0 {
		resetLetters() // Reset the letters if all have been used
	}

	rand.Seed(time.Now().UnixNano())
	index := rand.Intn(len(letters))
	letter := letters[index]
	letters = append(letters[:index], letters[index+1:]...)
	return letter
}

// getCurrentLetter returns the last letter selected if any remain.
func getCurrentLetter() rune {
	if len(letters) > 0 {
		return letters[len(letters)-1]
	}
	return ' ' // Return a space if no letters are left
}

// PetitBacHandler handles requests to the main game page, serving the game interface or handling form submission.
func PetitBacHandler(w http.ResponseWriter, r *http.Request) {
	scoreCookie, err := r.Cookie("score")
	currentScore := 0
	if err == nil {
		currentScore, _ = strconv.Atoi(scoreCookie.Value) // Convert the score from cookie if present
	}

	if r.Method == "POST" {
		// Handle form submission and save the response
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

	// Serve the game page with a new letter
	currentLetter := getRandomLetter()
	data := map[string]interface{}{
		"Letter": string(currentLetter),
		"Score":  currentScore,
	}
	if err := templates.ExecuteTemplate(w, "petit-bac.html", data); err != nil {
		http.Error(w, "Error rendering template: "+err.Error(), http.StatusInternalServerError)
	}
}

// AnswersHandler serves the page that displays the answers submitted by the user.
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

// ValidateAnswersHandler processes the form submission for validating answers and updating the score.
func ValidateAnswersHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	// Read and update the score based on correct answers
	scoreCookie, err := r.Cookie("score")
	currentScore := 0
	if err == nil {
		currentScore, _ = strconv.Atoi(scoreCookie.Value)
	}

	sessionScore := 0
	if r.FormValue("artistCorrect") == "true" {
		sessionScore++
	}
	if r.FormValue("albumCorrect") == "true" {
		sessionScore++
	}
	if r.FormValue("musicGroupCorrect") == "true" {
		sessionScore++
	}
	if r.FormValue("instrumentCorrect") == "true" {
		sessionScore++
	}
	if r.FormValue("featuringCorrect") == "true" {
		sessionScore++
	}

	currentScore += sessionScore

	// Save the new score in a cookie
	http.SetCookie(w, &http.Cookie{
		Name:   "score",
		Value:  strconv.Itoa(currentScore),
		Path:   "/",
		MaxAge: 86400, // The score expires in one day
	})

	// Display the updated score
	fmt.Fprintf(w, `<html><body><p>Your total score: %d</p><a href='/petit-bac'>Try again</a></body></html>`, currentScore)
}
