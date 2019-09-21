package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
)

var activeGames = make(map[string]*Game)

func createRouter() *mux.Router {
	router := mux.NewRouter()
	router.HandleFunc("/", serveHome)
	router.HandleFunc("/new", serveLobby)
	router.HandleFunc("/join", serveLobby)
	router.HandleFunc("/single-player", createSinglePlayerGame)
	router.HandleFunc("/g/{gameID}", serveGame)
	router.HandleFunc("/ws/{gameID}", func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		key := vars["gameID"]

		game := activeGames[key]
		if game != nil && game.available {
			connectPlayer(game, w, r)
			return
		}
	})

	router.HandleFunc("/api/new", createGame)
	router.HandleFunc("/api/join", joinGame)
	return router
}

func serveHome(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.Error(w, "Not found", http.StatusNotFound)
		return
	}
	if r.Method != "GET" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	http.ServeFile(w, r, "./frontend/html/home.html")
}

func createGame(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/api/new" {
		http.Error(w, "Not found", http.StatusNotFound)
		return
	}
	if r.Method != "GET" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	gameID, err := createGameAndWait()
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnprocessableEntity)
		return
	}
	response := make(map[string]interface{})
	response["gameId"] = gameID
	json.NewEncoder(w).Encode(response)
}

func createSinglePlayerGame(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/single-player" {
		http.Error(w, "Not found", http.StatusNotFound)
		return
	}
	if r.Method != "GET" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	gameID := randToken()
	game := newGame(gameID, height, width)
	activeGames[gameID] = game
	go game.run()

	connectBot(game)

	http.Redirect(w, r, fmt.Sprintf("/g/%s", gameID), http.StatusSeeOther)
}

func serveLobby(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/join" && r.URL.Path != "/new" {
		http.Error(w, "Not found", http.StatusNotFound)
		return
	}
	if r.Method != "GET" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	http.ServeFile(w, r, "./frontend/html/lobby.html")
}

func joinGame(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/api/join" {
		http.Error(w, "Not found", http.StatusNotFound)
		return
	}
	if r.Method != "GET" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	gameID, err := findAvailableGameAndJoin()
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnprocessableEntity)
		return
	}
	response := make(map[string]interface{})
	response["gameId"] = gameID
	json.NewEncoder(w).Encode(response)
}

func serveGame(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	vars := mux.Vars(r)
	key := vars["gameID"]

	game := activeGames[key]
	if game == nil || !game.available {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}
	http.ServeFile(w, r, "./frontend/html/game.html")
}
