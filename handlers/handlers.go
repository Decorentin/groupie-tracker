package handlers

import (
	"crypto/sha256"
	"database/sql"
	"encoding/hex"
	"fmt"
	"log"
	"math"
	"net/http"
	"test/user"
	"unicode"
)

// serveFile retourne un gestionnaire qui sert un fichier spécifique.
func serveFile(filename string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, filename)
	}
}

// LoginHandler gère les requêtes de connexion.
func LoginHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "POST" {
			login := r.FormValue("pseudo")
			password := r.FormValue("password")

			// Hachage du mot de passe avec SHA-256 avant de le comparer
			hashedPassword := hashPassword(password)

			loggedIn, err := user.LoginUser(db, login, hashedPassword)
			if err != nil {
				log.Println("Login failed:", err)
				http.Error(w, "Erreur de connexion", http.StatusInternalServerError)
				return
			}

			if loggedIn {
				http.Redirect(w, r, "/home", http.StatusSeeOther)
			} else {
				fmt.Fprintln(w, "Nom d'utilisateur ou mot de passe incorrect")
			}
		} else {
			serveFile("login.html")(w, r)
		}
	}
}

// RegisterHandler gère les requêtes d'inscription.
func RegisterHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "GET" {
			serveFile("register.html")(w, r)
		} else if r.Method == "POST" {
			pseudo := r.FormValue("pseudo")
			email := r.FormValue("email")
			password := r.FormValue("password")
			confirmPassword := r.FormValue("confirm_password")

			if password != confirmPassword {
				http.Error(w, "Les mots de passe ne correspondent pas", http.StatusBadRequest)
				return
			}

			// Vérification de la complexité du mot de passe
			passwordEntropy := passwordEntropy(password)
			if passwordEntropy < 80 { // Changement ici pour 80 bits
				http.Error(w, "Le mot de passe doit avoir une entropie d'au moins 80", http.StatusBadRequest)
				return
			}

			// Hachage du mot de passe avec SHA-256 avant de l'enregistrer
			hashedPassword := hashPassword(password)

			newUser := user.User{Pseudo: pseudo, Email: email, Password: hashedPassword}
			err := user.RegisterUser(db, newUser)
			if err != nil {
				log.Println("Failed to register user:", err)
				http.Error(w, "Failed to register user: pseudo or email already exists", http.StatusInternalServerError)
				return
			}

			http.Redirect(w, r, "/", http.StatusSeeOther)
		} else {
			http.Error(w, "Méthode non autorisée", http.StatusMethodNotAllowed)
		}
	}
}

// Fonction pour hasher un mot de passe avec l'algorithme SHA-256
func hashPassword(password string) string {
	hash := sha256.New()
	hash.Write([]byte(password))
	hashedPassword := hex.EncodeToString(hash.Sum(nil))
	return hashedPassword
}

// Fonction pour calculer l'entropie du mot de passe
func passwordEntropy(password string) int {
	var (
		entropy float64
		charset float64
	)

	for _, char := range password {
		if unicode.IsLetter(char) {
			charset += 52 // 26 lowercase + 26 uppercase
		} else if unicode.IsDigit(char) {
			charset += 10 // 10 digits
		} else {
			charset += 33 // 33 special characters
		}
	}

	entropy = float64(len(password)) * (math.Log2(charset))
	return int(entropy)
}
