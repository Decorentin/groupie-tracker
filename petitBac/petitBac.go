package petitBac

import (
	"encoding/json"
	"log"
	"math/rand"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
)

type Client struct {
	socket *websocket.Conn
	send   chan []byte
}

type Game struct {
	CurrentLetter rune
	Categories    map[string]string
}

var (
	clients    = make(map[*Client]bool)
	broadcast  = make(chan []byte)
	register   = make(chan *Client)
	unregister = make(chan *Client)
	upgrader   = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
	}
	currentGame *Game
	categories  = []string{"Artiste", "Album", "Groupe de musique", "Instrument de musique", "Featuring"} // Définition des catégories
)

func PetitBacHandler() {
	http.HandleFunc("/", serveHome)
	http.HandleFunc("/ws", handleWebSocket)

	initGame()
	go handleMessages()

	http.ListenAndServe(":8080", nil)
}

func serveHome(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "index.html")
}

func handleWebSocket(w http.ResponseWriter, r *http.Request) {
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("WebSocket upgrade failed: %v", err)
		return
	}
	client := &Client{socket: ws, send: make(chan []byte, 1024)}
	register <- client

	initMessage, _ := json.Marshal(map[string]interface{}{
		"letter":     string(currentGame.CurrentLetter),
		"categories": categories,
	})
	client.send <- initMessage

	go client.write()
	go client.read()
}

func initGame() {
	currentGame = &Game{
		CurrentLetter: generateUniqueLetter(),
		Categories:    make(map[string]string),
	}
}

func generateUniqueLetter() rune {
	rand.Seed(time.Now().UnixNano())
	return rune(rand.Intn(26) + 'A')
}

func (c *Client) read() {
	defer func() {
		unregister <- c
		c.socket.Close()
	}()

	for {
		_, message, err := c.socket.ReadMessage()
		if err != nil {
			log.Printf("error: %v", err)
			break
		}
		broadcast <- message
	}
}

func (c *Client) write() {
	defer c.socket.Close()
	for message := range c.send {
		err := c.socket.WriteMessage(websocket.TextMessage, message)
		if err != nil {
			log.Printf("error: %v", err)
			break
		}
	}
}

func handleMessages() {
	for {
		select {
		case client := <-register:
			clients[client] = true
		case client := <-unregister:
			if _, ok := clients[client]; ok {
				delete(clients, client)
				close(client.send)
			}
		case message := <-broadcast:
			// Afficher le message dans les logs du serveur
			log.Printf("Received message: %s", string(message))

			// Traitement supplémentaire des données (ex : stockage ou analyse)
			processData(message)

			// Diffuser le message à tous les clients
			for client := range clients {
				select {
				case client.send <- message:
				default:
					close(client.send)
					delete(clients, client)
				}
			}
		}
	}
}

func processData(message []byte) {
	// Deserialiser le message pour un traitement plus spécifique si nécessaire
	var data map[string]interface{}
	if err := json.Unmarshal(message, &data); err != nil {
		log.Printf("Error unmarshalling message: %v", err)
		return
	}

	// Traiter les données comme nécessaire
	log.Printf("Processed data: %+v", data)
}
