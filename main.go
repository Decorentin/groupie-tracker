package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"test/user"

	_ "github.com/mattn/go-sqlite3"
)

var db *sql.DB

func main() {

	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

	var err error
	db, err = sql.Open("sqlite3", "./foo.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// Créer la table 'users'
	_, err = db.Exec(`
CREATE TABLE IF NOT EXISTS users (
	id INTEGER PRIMARY KEY AUTOINCREMENT,
	username TEXT UNIQUE NOT NULL,
	email TEXT UNIQUE NOT NULL,
	password TEXT NOT NULL
)
`)
	if err != nil {
		log.Fatal(err)
	}

	// Insérer des données d'exemple dans la table 'users'
	// 	_, err = db.Exec(`
	// INSERT INTO users (username, email, password) VALUES ("testUser", "test@example.com", "testPassword")
	// `)
	// 	if err != nil {
	// 		log.Fatal(err)
	// 	}

	http.HandleFunc("/", loginHandler)
	http.HandleFunc("/login", loginHandler)
	http.HandleFunc("/home", serveFile("home.html"))
	http.HandleFunc("/register", registerHandler)
	http.HandleFunc("/signin", registerHandler)

	http.ListenAndServe(":8080", nil)
}

func loginHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		// Récupérer les données du formulaire
		login := r.FormValue("username")
		password := r.FormValue("password")

		// Vérifier si les informations d'identification sont correctes
		loggedIn, err := user.LoginUser(db, login, password)
		if err != nil {
			log.Println("Login failed:", err)
			http.Error(w, "Erreur de connexion", http.StatusInternalServerError)
			return
		}

		if loggedIn {
			// Rediriger l'utilisateur vers la page 'home'
			http.Redirect(w, r, "/home", http.StatusSeeOther)
		} else {
			// Afficher un message d'erreur
			fmt.Fprintln(w, "Nom d'utilisateur ou mot de passe incorrect")
		}
	} else {
		serveFile("login.html")(w, r)
	}
}

func registerHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		serveFile("register.html")(w, r)
	} else if r.Method == "POST" {
		// Récupérer les données du formulaire
		username := r.FormValue("username")
		email := r.FormValue("email")
		password := r.FormValue("password")

		// Enregistrer l'utilisateur dans la base de données
		newUser := user.User{Username: username, Email: email, Password: password}
		err := user.RegisterUser(db, newUser)
		if err != nil {
			log.Println("Failed to register user:", err)
			http.Error(w, "Failed to register user: username or email already exists", http.StatusInternalServerError)
			return
		}

		// Rediriger l'utilisateur vers la page de connexion
		http.Redirect(w, r, "/", http.StatusSeeOther)
	} else {
		http.Error(w, "Méthode non autorisée", http.StatusMethodNotAllowed)
	}
}

func serveFile(filename string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, filename)
	}
}
