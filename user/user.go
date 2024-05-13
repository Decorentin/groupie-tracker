package user

import (
	"database/sql"
	"fmt"
)

// User represents a platform user
type User struct {
	Pseudo   string
	Email    string
	Password string
}

// RegisterUser registers a new user
func RegisterUser(db *sql.DB, user User) error {
	var exists int
	err := db.QueryRow("SELECT COUNT(*) FROM USER WHERE pseudo = ? OR email = ?", user.Pseudo, user.Email).Scan(&exists)
	if err != nil {
		return err
	}
	if exists > 0 {
		return fmt.Errorf("pseudo or email already exists")
	}

	_, err = db.Exec("INSERT INTO USER (pseudo, email, password) VALUES (?, ?, ?)", user.Pseudo, user.Email, user.Password)
	return err
}

// LoginUser allows a user to log in
func LoginUser(db *sql.DB, login, password string) (bool, error) {
	var dbPassword string
	var err error

	// Search for the user using the pseudo
	err = db.QueryRow("SELECT password FROM USER WHERE pseudo = ?", login).Scan(&dbPassword)
	if err != nil {
		// If the user is not found using the pseudo, search using the email address
		err = db.QueryRow("SELECT password FROM USER WHERE email = ?", login).Scan(&dbPassword)
		if err != nil {
			if err == sql.ErrNoRows {
				return false, fmt.Errorf("user not found")
			}
			return false, err
		}
	}

	// Check the password
	if dbPassword != password {
		return false, fmt.Errorf("incorrect password")
	}

	return true, nil
}
