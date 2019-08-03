package main

import (
	"fmt"
	"math/rand"
	"sync"
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
