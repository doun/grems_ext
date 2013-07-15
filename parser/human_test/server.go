package main

import (
	"log"
	"net/http"
)

func main() {

	http.HandleFunc("/parser", parser)
	err := http.ListenAndServe(":909", nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}

func parser(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Access-Control-Allow-Origin", "*")
	w.Header().Add("Access-Control-Allow-Headers", "Content-Type")
	_, h, e := r.FormFile("file")
	if e != nil {
		w.Write([]byte("no..."))
	} else {
		w.Write([]byte("ok!" + h.Filename))
	}
}
