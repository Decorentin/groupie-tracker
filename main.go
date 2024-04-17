package main

import (
	"database/sql"
	"log"
	"net/http"
	guessthesong "test/guessTheSong"
	"test/handlers"

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
CREATE TABLE IF NOT EXISTS USER (
	id INTEGER PRIMARY KEY,
	username TEXT UNIQUE NOT NULL,
	email TEXT UNIQUE NOT NULL,
	password TEXT NOT NULL
);

CREATE TABLE IF NOT EXISTS ROOMS (
    id INTEGER PRIMARY KEY,
    created_by INTEGER NOT NULL,
    max_player INTEGER NOT NULL,
    name TEXT NOT NULL,
    id_game INTEGER,
    FOREIGN KEY (created_by) REFERENCES USER(id),
    FOREIGN KEY (id_game) REFERENCES GAMES(id)
);

CREATE TABLE IF NOT EXISTS ROOM_USERS (
    id_room INTEGER,
    id_user INTEGER,
    score INTEGER,
    FOREIGN KEY (id_room) REFERENCES ROOMS(id),
    FOREIGN KEY (id_user) REFERENCES USER(id),
    PRIMARY KEY (id_room, id_user)
);

CREATE TABLE IF NOT EXISTS GAMES (
    id INTEGER PRIMARY KEY,
    name TEXT NOT NULL
)


`)
	if err != nil {
		log.Fatal(err)
	}

	http.HandleFunc("/", handlers.LoginHandler(db))
	http.HandleFunc("/login", handlers.LoginHandler(db))
	http.HandleFunc("/register", handlers.RegisterHandler(db))

	http.HandleFunc("/home", serveFile("home.html"))

	http.HandleFunc("/guess-the-song", guessthesong.GuessTheSongHandler)

	http.ListenAndServe(":8080", nil)
}

// serveFile retourne un gestionnaire qui sert un fichier spécifique.
func serveFile(filename string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, filename)
	}
}
