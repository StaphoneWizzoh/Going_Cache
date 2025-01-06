package main

import (
	"fmt"
	"net/http"
)

func main() {
	cs := NewCacheServer()
	http.HandleFunc("/set", cs.SetHandler)
	http.HandleFunc("/get", cs.GetHandler)
	err := http.ListenAndServe(":8080", nil)
	if err != nil{
		fmt.Println(err)
		return
	}
}
