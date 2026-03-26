package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"strings"
)

type Photo struct {
	Src     string `json:"src"`
	Caption string `json:"caption"`
}

type Species struct {
	Slug       string  `json:"slug"`
	Group      string  `json:"group"`
	Name       string  `json:"name"`
	Scientific string  `json:"scientific"`
	Info       string  `json:"info"`
	Habitat    string  `json:"habitat"`
	Breeding   string  `json:"breeding"`
	Captive    string  `json:"captive"`
	Redlist    string  `json:"redlist"`
	Photos     []Photo `json:"photos"`
}

var speciesData []Species

func loadSpecies() {
	file, err := os.Open("api/species.json")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()
	json.NewDecoder(file).Decode(&speciesData)
}

func apiAllSpecies(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(speciesData)
}

func apiSingleSpecies(w http.ResponseWriter, r *http.Request) {
	slug := strings.TrimPrefix(r.URL.Path, "/api/species/")
	for _, s := range speciesData {
		if s.Slug == slug {
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(s)
			return
		}
	}
	http.NotFound(w, r)
}

func main() {
	loadSpecies()

	// API routes first
	http.HandleFunc("/api/species/", apiSingleSpecies)
	http.HandleFunc("/api/species", apiAllSpecies)

	// Serve static files
	http.Handle("/js/", http.StripPrefix("/js/", http.FileServer(http.Dir("static/js"))))
	http.Handle("/css/", http.StripPrefix("/css/", http.FileServer(http.Dir("static/css"))))
	http.Handle("/images/", http.StripPrefix("/images/", http.FileServer(http.Dir("static/images"))))

	// SPA entry point
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "static/index.html")
	})

	log.Println("Server running at :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
