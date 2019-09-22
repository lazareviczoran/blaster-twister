package main

import (
	"flag"
	"log"
	"math/rand"
	"net/http"
	"os"
	"time"
)

func main() {
	HOST := ""
	if os.Getenv("GO_ENV") == "development" {
		HOST = "localhost"
	}
	PORT := os.Getenv("PORT")
	if PORT == "" {
		PORT = "8080"
	}
	rand.Seed(time.Now().UnixNano())

	flag.Parse()
	router := createRouter()
	http.Handle("/src/", http.StripPrefix("/src/", http.FileServer(http.Dir("./dist"))))
	http.Handle("/css/", http.StripPrefix("/css/", http.FileServer(http.Dir("./frontend/css"))))
	http.Handle("/img/", http.StripPrefix("/img/", http.FileServer(http.Dir("./frontend/img"))))
	http.Handle("/", router)
	err := http.ListenAndServe(HOST+":"+PORT, nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
