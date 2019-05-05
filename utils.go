package main

import (
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
