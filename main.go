package main

import (
	"fmt"
	"log"
	"net/http"
	"time"
)

func main() {
	cs := NewCacheServer()
	cache := NewCache(5)
	cache.startEvictionTicker(1 * time.Minute)

	http.HandleFunc("/set", cs.SetHandler)
	http.HandleFunc("/get", cs.GetHandler)
	
	log.Println("Server listening at port: 8080")
	err := http.ListenAndServe(":8080", nil)
	if err != nil{
		fmt.Println(err)
		return
	}
}
