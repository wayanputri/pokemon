package main

import (
	"embed"
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"
	"strings"
)

//go:embed views/*
var views embed.FS
var t = template.Must(template.ParseFS(views, "views/*"))

func main() {
	router := http.NewServeMux()
	router.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		resp, err := http.Get("https://pokeapi.co/api/v2/pokemon/pikachu")
		if err != nil {
			http.Error(w, "Unable to grab the pokemon data", http.StatusInternalServerError)
		}
		data := Pokemon{}
		if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
			http.Error(w, "Unable to parse the pokemon data", http.StatusInternalServerError)
		}
		if err := t.ExecuteTemplate(w, "index.html", data); err != nil {
			http.Error(w, "Something went wrong", http.StatusInternalServerError)
		}

	})
	router.HandleFunc("/poke", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
			return
		}

		if err := r.ParseForm(); err != nil {
			http.Error(w, "Unable to parse form", http.StatusInternalServerError)
			return
		}
		formValue := strings.ToLower(r.FormValue("pokemon"))

		resp, err := http.Get("https://pokeapi.co/api/v2/pokemon/" + formValue)
		if err != nil {
			http.Error(w, "Unable to fetch the pokemon data", http.StatusInternalServerError)
		}
		data := Pokemon{}
		if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
			http.Error(w, "Unable to parse the pokemon data", http.StatusInternalServerError)
		}
		if err := t.ExecuteTemplate(w, "response.html", data); err != nil {
			http.Error(w, "Something went wrong", http.StatusInternalServerError)
		}
		fmt.Println("data post", data)

	})
	server := http.Server{
		Addr:    ":3000",
		Handler: router,
	}

	fmt.Println("Listening of :3000")
	server.ListenAndServe()
}
