package user

import (
	"database/sql"
	"fmt"
)

// User représente un utilisateur de la plateforme
type User struct {
	Username string
	Email    string
	Password string
}

// RegisterUser permet d'inscrire un nouvel utilisateur
func RegisterUser(db *sql.DB, user User) error {
	var exists int
	err := db.QueryRow("SELECT COUNT(*) FROM users WHERE username = ? OR email = ?", user.Username, user.Email).Scan(&exists)
	if err != nil {
		return err
	}
	if exists > 0 {
		return fmt.Errorf("username or email already exists")
	}

	_, err = db.Exec("INSERT INTO users (username, email, password) VALUES (?, ?, ?)", user.Username, user.Email, user.Password)
	return err
}

// LoginUser permet à un utilisateur de se connecter
func LoginUser(db *sql.DB, login, password string) (bool, error) {
	var dbPassword string
	var err error

	// Rechercher l'utilisateur en utilisant le pseudo
	err = db.QueryRow("SELECT password FROM users WHERE username = ?", login).Scan(&dbPassword)
	if err != nil {
		// Si l'utilisateur n'est pas trouvé en utilisant le pseudo, rechercher en utilisant l'adresse e-mail
		err = db.QueryRow("SELECT password FROM users WHERE email = ?", login).Scan(&dbPassword)
		if err != nil {
			if err == sql.ErrNoRows {
				return false, fmt.Errorf("utilisateur non trouvé")
			}
			return false, err
		}
	}

	// Vérifier le mot de passe
	if dbPassword != password {
		return false, fmt.Errorf("mot de passe incorrect")
	}

	return true, nil
}

