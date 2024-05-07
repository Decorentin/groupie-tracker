package dbHandlers

import (
	"crypto/sha256"
	"database/sql"
	"encoding/hex"
	"encoding/json"
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
			charset += 52
		} else if unicode.IsDigit(char) {
			charset += 10
		} else {
			charset += 33
		}
	}

	entropy = float64(len(password)) * (math.Log2(charset))
	return int(entropy)
}

// Room représente une salle dans l'application.
type Room struct {
	ID         int    `json:"id"`
	CreatedBy  int    `json:"created_by"`
	MaxPlayers int    `json:"max_players"`
	Name       string `json:"name"`
	GameID     int    `json:"game_id"`
}

// CreateRoomHandler gère la création d'une nouvelle salle dans la base de données.
func CreateRoomHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var room Room
		err := json.NewDecoder(r.Body).Decode(&room)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		_, err = db.Exec("INSERT INTO ROOMS (created_by, max_player, name, id_game) VALUES (?, ?, ?, ?)",
			room.CreatedBy, room.MaxPlayers, room.Name, room.GameID)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusCreated)
	}
}

// ListRoomsHandler récupère la liste des salles disponibles à partir de la base de données.
func ListRoomsHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		rows, err := db.Query("SELECT id, created_by, max_player, name, id_game FROM ROOMS")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		defer rows.Close()

		var rooms []Room
		for rows.Next() {
			var room Room
			err := rows.Scan(&room.ID, &room.CreatedBy, &room.MaxPlayers, &room.Name, &room.GameID)
			if err != nil {
				log.Println(err)
				continue
			}
			rooms = append(rooms, room)
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(rooms)
	}
}

// JoinRoomHandler permet aux utilisateurs de rejoindre une salle existante.
func JoinRoomHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var requestData struct {
			RoomID int `json:"room_id"`
			UserID int `json:"user_id"`
		}
		err := json.NewDecoder(r.Body).Decode(&requestData)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		var roomExists bool
		err = db.QueryRow("SELECT EXISTS (SELECT 1 FROM ROOMS WHERE id = ?)", requestData.RoomID).Scan(&roomExists)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		if !roomExists {
			http.Error(w, "La salle spécifiée n'existe pas", http.StatusBadRequest)
			return
		}

		var userExists bool
		err = db.QueryRow("SELECT EXISTS (SELECT 1 FROM USER WHERE id = ?)", requestData.UserID).Scan(&userExists)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		if !userExists {
			http.Error(w, "L'utilisateur spécifié n'existe pas", http.StatusBadRequest)
			return
		}

		// Logique pour rejoindre la salle
		// Par exemple, vous pourriez insérer une entrée dans la table ROOM_USERS

		w.WriteHeader(http.StatusOK)
	}
}

// GetRoomDetailsHandler récupère les détails d'une salle spécifique à partir de son ID.
func GetRoomDetailsHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		roomID := r.URL.Query().Get("id")

		var room Room
		err := db.QueryRow("SELECT id, created_by, max_player, name, id_game FROM ROOMS WHERE id = ?", roomID).
			Scan(&room.ID, &room.CreatedBy, &room.MaxPlayers, &room.Name, &room.GameID)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(room)
	}
}

// DeleteRoomHandler supprime une salle existante de la base de données.
func DeleteRoomHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		roomID := r.URL.Query().Get("id")

		_, err := db.Exec("DELETE FROM ROOMS WHERE id = ?", roomID)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
	}
}
