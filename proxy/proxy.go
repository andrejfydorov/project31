package main

import (
	"bytes"
	"io/ioutil"
	"log"
	"net/http"
)

const proxyAddr string = "localhost:9000"

var (
	counter            int    = 0
	firstInstanceHost  string = "http://localhost:8080"
	secondInstanceHost string = "http://localhost:8081"
)

func main() {
	http.HandleFunc("/", handleProxy)
	http.HandleFunc("/create", handleProxy)
	http.HandleFunc("/make_friends", handleProxy)
	http.HandleFunc("/user", handleProxy)
	http.HandleFunc("/friends/{id:[0-9]+}", handleProxy)
	http.HandleFunc("/{id:[0-9]+}", handleProxy)
	http.HandleFunc("/get", handleProxy)

	log.Fatalln(http.ListenAndServe(proxyAddr, nil))
}

func handleProxy(w http.ResponseWriter, r *http.Request) {

	content, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}
	defer r.Body.Close()

	if counter == 0 {
		resp, err := http.Post(firstInstanceHost+r.URL.Path, "text/json", bytes.NewBuffer(content))
		if err != nil {
			log.Fatalln(err)
		}
		respBytes, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			log.Fatalln(err)
		}
		defer resp.Body.Close()

		w.WriteHeader(http.StatusOK)
		w.Write(respBytes)

		counter++
		return
	}

	resp, err := http.Post(firstInstanceHost+r.URL.Path, "text/json", bytes.NewBuffer(content))
	if err != nil {
		log.Fatalln(err)
	}
	respBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalln(err)
	}
	defer resp.Body.Close()

	w.WriteHeader(http.StatusOK)
	w.Write(respBytes)

	counter--
	return

}
