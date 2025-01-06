package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"
)

var port string
var peers string

func main() {

	flag.StringVar(&port, "port", ":8080", "HTTP server port")
	flag.StringVar(&peers, "peers", "", "Comma-separated list of peer addresses")
	flag.Parse()

	peerList := strings.Split(peers, ",")
	cs := NewCacheServer(peerList)
	cache := NewCache(10)
	cache.startEvictionTicker(1 * time.Minute)

	http.HandleFunc("/set", cs.SetHandler)
	http.HandleFunc("/get", cs.GetHandler)

	log.Println("Server listening at port: 8080")
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		fmt.Println(err)
		return
	}
}
