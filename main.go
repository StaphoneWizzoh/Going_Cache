package main

import (
	"flag"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"strings"
)

var port string
var peers string

func main() {

	flag.StringVar(&port, "port", ":8080", "HTTP server port")
	flag.StringVar(&peers, "peers", "", "Comma-separated list of peer addresses")
	flag.Parse()

	nodeID := fmt.Sprintf("%s%d", "node", rand.Intn(100))
	peerList := strings.Split(peers, ",")
	cs := NewCacheServer(peerList, nodeID)

	http.HandleFunc("/set", cs.SetHandler)
	http.HandleFunc("/get", cs.GetHandler)

	log.Printf("Server listening at port: %s\n", port)
	log.Fatal(http.ListenAndServe(port, nil))

}
