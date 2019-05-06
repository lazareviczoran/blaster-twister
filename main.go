package main

import (
	"flag"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"time"

	"github.com/gorilla/mux"
)

const height int = 400
const width int = 300

var addr = flag.String("addr", "localhost:8080", "http service address")
var activeGames = make(map[string]*Game)

func serveHome(w http.ResponseWriter, r *http.Request) {
	log.Printf("serveHome")
	if r.URL.Path != "/" {
		http.Error(w, "Not found", http.StatusNotFound)
		return
	}
	if r.Method != "GET" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	http.ServeFile(w, r, "home.html")
}

func createGame(w http.ResponseWriter, r *http.Request) {
	log.Printf("createGame")
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

func serveGame(w http.ResponseWriter, r *http.Request) {
	log.Printf("serveGame")
	if r.Method != "GET" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	http.ServeFile(w, r, "game.html")
}

func main() {
	rand.Seed(time.Now().UnixNano())

	flag.Parse()

	router := mux.NewRouter()
	router.HandleFunc("/", serveHome)
	router.HandleFunc("/new", createGame)
	router.HandleFunc("/g/{gameID}", serveGame)
	router.HandleFunc("/ws/{gameID}", func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		key := vars["gameID"]

		game := activeGames[key]
		if game != nil {
			connect(game, w, r)
		}
	})
	http.Handle("/src/", http.StripPrefix("/src/", http.FileServer(http.Dir("./scripts"))))
	http.Handle("/", router)
	err := http.ListenAndServe(*addr, nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}

func getStartRotation() int {
	return rand.Intn(360)
}

func randToken() string {
	b := make([]byte, 8)
	rand.Read(b)
	return fmt.Sprintf("%x", b)
}
