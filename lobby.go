package main

import (
	"errors"
	"log"
)

var lobby *Lobby

type Lobby struct {
	activeCandidates []Candidate
	register         chan Candidate
}

func initLobby() {
	if lobby == nil {
		lobby = &Lobby{
			make([]Candidate, 0),
			make(chan Candidate),
		}
	}
	go lobby.run()
}

func (l *Lobby) run() {
	for {
		select {
		case candidate := <-l.register:
			{
				l.activeCandidates = append(l.activeCandidates, candidate)
				go l.tryToStart()
			}
		}
	}
}

func (l *Lobby) tryToStart() {
	if len(l.activeCandidates) >= 2 {
		cand1, err := l.getReadyCandidate()
		if err != nil {
			return
		}
		cand2, err := l.getReadyCandidate()
		if err != nil {
			return
		}

		gameID, err := createGame()
		if err != nil {
			log.Printf("Error while starting game %s", err.Error())
		}
		cand1.Redirect([]byte(gameID))
		cand2.Redirect([]byte(gameID))
	}
}

func (l *Lobby) getReadyCandidate() (*Candidate, error) {
	found := false
	var cand Candidate
	for !found {
		if len(l.activeCandidates) == 0 {
			return nil, errors.New("There are no active players to join")
		}
		cand = l.activeCandidates[0]
		l.activeCandidates = l.activeCandidates[1:]
		found = cand.IsConnected()
	}
	return &cand, nil
}
