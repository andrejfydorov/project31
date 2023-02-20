package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
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

	if counter == 0 {

		originServerURL, err := url.Parse(firstInstanceHost)
		if err != nil {
			log.Fatal("invalid origin server URL")
		}

		r.Host = originServerURL.Host
		r.URL.Host = originServerURL.Host

		originServerResponse, err := http.DefaultClient.Do(r)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			_, _ = fmt.Fprint(w, err)
			return
		}

		w.WriteHeader(http.StatusOK)
		io.Copy(w, originServerResponse.Body)

		counter++
		return
	}

	originServerURL, err := url.Parse(secondInstanceHost)
	if err != nil {
		log.Fatal("invalid origin server URL")
	}

	r.Host = originServerURL.Host
	r.URL.Host = originServerURL.Host

	originServerResponse, err := http.DefaultClient.Do(r)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = fmt.Fprint(w, err)
		return
	}

	w.WriteHeader(http.StatusOK)
	io.Copy(w, originServerResponse.Body)

	counter--
	return

}
