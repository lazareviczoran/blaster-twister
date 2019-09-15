package main

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
)

var activeGames = make(map[string]*Game)

func createRouter() *mux.Router {
	router := mux.NewRouter()
	router.HandleFunc("/", serveHome)
	router.HandleFunc("/new", createGame)
	router.HandleFunc("/join", joinGame)
	router.HandleFunc("/single-player", createSinglePlayerGame)
	router.HandleFunc("/g/{gameID}", serveGame)
	router.HandleFunc("/ws/{gameID}", func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		key := vars["gameID"]

		game := activeGames[key]
		if game != nil {
			connectPlayer(game, w, r)
		}
	})
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
	if r.URL.Path != "/new" {
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

	http.Redirect(w, r, fmt.Sprintf("/g/%s", gameID), http.StatusSeeOther)
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

func joinGame(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/join" {
		http.Error(w, "Not found", http.StatusNotFound)
		return
	}
	if r.Method != "GET" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var oldestGame *Game
	for _, game := range activeGames {
		if game.available && (oldestGame == nil || oldestGame.createdAt.After(game.createdAt)) {
			oldestGame = game
		}
	}
	if oldestGame == nil {
		http.ServeFile(w, r, "./frontend/html/unavailable.html")
		return
	}
	http.Redirect(w, r, fmt.Sprintf("/g/%s", oldestGame.id), http.StatusSeeOther)
}

func serveGame(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	http.ServeFile(w, r, "./frontend/html/game.html")
}
