package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
)

type Link struct {
	ID        int
	ShortCode string `json:"short_code"`
	LongUrl   string `json:"long_url"`
}

var LinkList []Link

func Urlshort(w http.ResponseWriter, r *http.Request) {
	var reqLinkStruct Link

	err := json.NewDecoder(r.Body).Decode(&reqLinkStruct)
	if err != nil {
		http.Error(w, "Invalid Json", http.StatusBadRequest)
		return
	}

	reqLinkStruct.ID = len(LinkList) + 1

	reqLinkStruct.ShortCode = strconv.Itoa(reqLinkStruct.ID)

	LinkList = append(LinkList, reqLinkStruct)

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(reqLinkStruct)

}

func Redirect(w http.ResponseWriter, r *http.Request) {
	shortCode := r.PathValue("shortCode")

	if len(shortCode) == 0 {
		http.Error(w, "Invalid short code", http.StatusBadRequest)
		return
	}

	for _, link := range LinkList {
		if shortCode == link.ShortCode {
			http.Redirect(w, r, link.LongUrl, http.StatusFound)
			return
		}
	}

	http.Error(w, "Invalid Short Code", http.StatusBadRequest)

}

func main() {
	mux := http.NewServeMux()

	mux.Handle("POST /links", http.HandlerFunc(Urlshort))
	mux.Handle("GET /{shortCode}", http.HandlerFunc(Redirect))

	fmt.Println("Running server on port 3000")
	http.ListenAndServe(":3000", mux)

}
