package main

import (
	"net/http"
)

const (
	projectID = "be-matsuhashi-tatsuya"
	baseURL   = "https://azami-jp.com/live"
)

func main() {
	http.HandleFunc("/", handler)
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		return
	}
}

func handler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	client, err := newDatastoreClient(ctx, projectID)
	if err != nil {
		panic(err)
	}
	defer client.Close()

	lives, err := parse(baseURL)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = saveLatestLives(ctx, client, lives)
	if err != nil {
		panic(err)
	}

	err = saveLiveMaster(ctx, client, lives)
	if err != nil {
		panic(err)
	}
}
