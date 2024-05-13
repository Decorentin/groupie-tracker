package dbHandlers

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

// serveFile returns a handler that serves a specific file.
func serveFile(filename string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, filename)
	}
}

// LoginHandler handles login requests.
func LoginHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "POST" {
			login := r.FormValue("pseudo")
			password := r.FormValue("password")

			// Hash the password with SHA-256 before comparing it
			hashedPassword := hashPassword(password)

			loggedIn, err := user.LoginUser(db, login, hashedPassword)
			if err != nil {
				log.Println("Login failed:", err)
				http.Error(w, "Login error", http.StatusInternalServerError)
				return
			}

			if loggedIn {
				http.Redirect(w, r, "/home", http.StatusSeeOther)
			} else {
				fmt.Fprintln(w, "Incorrect username or password")
			}
		} else {
			serveFile("login.html")(w, r)
		}
	}
}

// RegisterHandler handles registration requests.
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
				http.Error(w, "Passwords do not match", http.StatusBadRequest)
				return
			}

			// Check the complexity of the password
			passwordEntropy := passwordEntropy(password)
			if passwordEntropy < 80 { // Change here to 80 bits
				http.Error(w, "Password must have an entropy of at least 80", http.StatusBadRequest)
				return
			}

			// Hash the password with SHA-256 before saving it
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
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	}
}

// Function to hash a password with the SHA-256 algorithm
func hashPassword(password string) string {
	hash := sha256.New()
	hash.Write([]byte(password))
	hashedPassword := hex.EncodeToString(hash.Sum(nil))
	return hashedPassword
}

// Function to calculate the entropy of the password
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
