package main

import (
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

var activeGames = make(map[string]*Game)

func createRouter() *mux.Router {
	router := mux.NewRouter()
	router.HandleFunc("/", serveHome)
	router.HandleFunc("/join", serveLobby)
	router.HandleFunc("/single-player", createSinglePlayerGame)
	router.HandleFunc("/g/{gameID}", serveGame)
	router.HandleFunc("/ws/game/{gameID}", func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		key := vars["gameID"]

		game := activeGames[key]
		if game != nil && !game.started {
			connectPlayer(game, w, r)
			return
		}
	})
	router.HandleFunc("/ws/game/{gameID}/{clientID}/{playerID}", func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		key := vars["gameID"]
		playerID, err := strconv.Atoi(vars["playerID"])
		if err != nil {
			log.Printf("Couldn't parse playerID %s", vars["playerID"])
			return
		}
		clientID, err := strconv.Atoi(vars["clientID"])
		if err != nil {
			log.Printf("Couldn't parse clientID %s", vars["clientID"])
			return
		}

		game := activeGames[key]
		if game != nil {
			player := game.players[playerID]
			if player != nil {
				if player.ClientID() != clientID {
					log.Printf("Client IDs don't match (%d != %d)", player.ClientID(), clientID)
					return
				}
				conn, err := upgrader.Upgrade(w, r, nil)
				if err != nil {
					log.Print("upgrade:", err)
					return
				}
				player.(*Human).AttachWriteConn(conn)
				return
			} else {
				log.Printf("Couldn't find player %d in game %s)", playerID, key)
			}
		}
		http.Error(w, "Not found", http.StatusNotFound)
	})
	router.HandleFunc("/ws/lobby", func(w http.ResponseWriter, r *http.Request) {
		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			log.Print("ws/lobby upgrade:", err)
			return
		}

		lobby.register <- newCandidate(conn)
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
	if r.URL.Path != "/join" {
		http.Error(w, "Not found", http.StatusNotFound)
		return
	}
	if r.Method != "GET" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	http.ServeFile(w, r, "./frontend/html/lobby.html")
}

func serveGame(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	vars := mux.Vars(r)
	key := vars["gameID"]

	game := activeGames[key]
	if game == nil || game.started {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}
	http.ServeFile(w, r, "./frontend/html/game.html")
}
