package main

import (
	"fmt"
	"net/http"
)

func main() {
	http.HandleFunc("/", handler)
	http.ListenAndServe(":8080", nil)
}

func handler(w http.ResponseWriter, r *http.Request) {
	baseURL := "https://azami-jp.com/live"

	// Goquery
	resp, err := fetch(baseURL)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	infosByGoquery, err := parseByGoquery(resp)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	fmt.Fprintln(w, "Goquery")
	for _, item := range infosByGoquery {
		fmt.Fprintln(w, item)
		// fmt.Fprintln(w, info, info.URL) // TODO ランタイムのバグか？出力が異なる
	}

	// Colly
	infosByColly, err := parseByColly(baseURL)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	fmt.Fprintln(w, "====================================")
	fmt.Fprintln(w, "Colly")
	for _, info := range infosByColly {
		fmt.Fprintln(w, info)
		// fmt.Fprintln(w, info, info.URL) // TODO ランタイムのバグか？出力が異なる
	}
}
