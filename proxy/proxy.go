package main

import (
	"bytes"
	"fmt"
	"io"
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

	log.Fatalln(http.ListenAndServe(proxyAddr, nil))
}

func handleProxy(w http.ResponseWriter, r *http.Request) {

	if counter == 0 {

		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		defer r.Body.Close()

		url := fmt.Sprintf("%v%v", firstInstanceHost, r.URL)

		newreq, err := http.NewRequest(r.Method, url, bytes.NewReader(body))
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		newreq.Header.Set("Content-Type", "application/json")

		client := &http.Client{}

		originServerResponse, err := client.Do(newreq)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
		io.Copy(w, originServerResponse.Body)

		counter++
		return
	}

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	url := fmt.Sprintf("%v%v", secondInstanceHost, r.URL)

	newreq, err := http.NewRequest(r.Method, url, bytes.NewReader(body))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	newreq.Header.Set("Content-Type", "application/json")

	client := &http.Client{}

	originServerResponse, err := client.Do(newreq)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	io.Copy(w, originServerResponse.Body)

	counter--
	return

}
