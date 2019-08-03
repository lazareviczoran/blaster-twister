package main

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
)

func syncMapToMap(m *sync.Map) map[string]interface{} {
	temp := make(map[string]interface{})
	m.Range(func(k interface{}, v interface{}) bool {
		temp[k.(string)] = v
		return true
	})

	return temp
}

func getStartRotation() int {
	return rand.Intn(90)
}

func randToken() string {
	b := make([]byte, 8)
	rand.Read(b)
	return fmt.Sprintf("%x", b)
}

func randomIntFromRange(min int, max int) int {
	return rand.Intn(max-min) + min
}

func createRandomIntervalTicker(min, max int) *time.Ticker {
	interval := randomIntFromRange(min, max)
	return time.NewTicker(time.Duration(interval) * time.Millisecond)
}
