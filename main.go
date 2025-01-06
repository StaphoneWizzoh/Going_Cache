package main

import (
	"fmt"
	"log"
	"net/http"
)

func main() {
	cs := NewCacheServer()
	http.HandleFunc("/set", cs.SetHandler)
	http.HandleFunc("/get", cs.GetHandler)
	log.Println("Server listening at port: 8080")
	err := http.ListenAndServe(":8080", nil)
	if err != nil{
		fmt.Println(err)
		return
	}
}
