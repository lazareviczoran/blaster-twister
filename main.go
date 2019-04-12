package main

import (
	"flag"
	"log"
	"math/rand"
	"net/http"
)

const height int = 400
const width int = 300

var addr = flag.String("addr", "localhost:8080", "http service address")

func serveHome(w http.ResponseWriter, r *http.Request) {
	log.Println(r.URL)
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

func main() {
	// rand.Seed(time.Now().UnixNano())
	// for i := 0; i < 2; i++ {
	// 	fmt.Printf("Player%d's staring rotation is:%d\n", i, getStartRotation())
	// }
	// board := initBoard(height, width)
	// board.fields[30][50].isUsed = true

	// fmt.Printf("This is the value of the 30th row and 50th column: %+v\n", &board.fields[30][50])

	flag.Parse()
	game := newGame()
	go game.run()

	http.HandleFunc("/", serveHome)
	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		connect(game, w, r)
	})
	err := http.ListenAndServe(*addr, nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}

func getStartRotation() int {
	return rand.Intn(360)
}
