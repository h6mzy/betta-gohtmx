package main

import (
	"encoding/json"
	"fmt"
	"html/template"
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

type PageData struct {
	Species        []Species
	CurrentSpecies *Species
}

var speciesList []Species
var tmpl *template.Template

func loadSpecies() {
	data, err := os.ReadFile("data/species.json")
	if err != nil {
		log.Fatal(err)
	}
	if err := json.Unmarshal(data, &speciesList); err != nil {
		log.Fatal(err)
	}
}

func homeHandler(w http.ResponseWriter, r *http.Request) {
	page := PageData{Species: speciesList}
	tmpl.ExecuteTemplate(w, "layout.html", page)
}

func speciesHandler(w http.ResponseWriter, r *http.Request) {
	slug := strings.TrimPrefix(r.URL.Path, "/species/")
	var current *Species
	for _, s := range speciesList {
		if s.Slug == slug {
			current = &s
			break
		}
	}
	if current == nil {
		http.NotFound(w, r)
		return
	}

	// HTMX partial request
	if r.Header.Get("HX-Request") == "true" {
		tmpl.ExecuteTemplate(w, "species_partial", current) // Use template **name**, not filename
		return
	}

	// Full page for direct URL
	page := PageData{
		Species:        speciesList,
		CurrentSpecies: current,
	}
	tmpl.ExecuteTemplate(w, "layout.html", page)
}

func apiSpeciesHandler(w http.ResponseWriter, r *http.Request) {
	query := strings.ToLower(r.URL.Query().Get("q"))

	w.Header().Set("Content-Type", "text/html")

	for _, s := range speciesList {
		if query == "" ||
			strings.Contains(strings.ToLower(s.Name), query) ||
			strings.Contains(strings.ToLower(s.Group), query) ||
			strings.Contains(strings.ToLower(s.Redlist), query) {
			fmt.Fprintf(w,
				`<div>
					<a hx-get="/species/%s"
					   hx-target="#species-detail"
					   hx-swap="innerHTML"
					   hx-push-url="true">
					   %s
					</a> — %s — %s
				</div>`,
				s.Slug, s.Name, s.Group, s.Redlist)
		}
	}
}

func main() {
	loadSpecies()
	tmpl = template.Must(template.ParseFiles(
		"templates/layout.html",
		"templates/home.html",
		"templates/species_partial.html",
	))

	http.HandleFunc("/", homeHandler)
	http.HandleFunc("/species/", speciesHandler)
	http.HandleFunc("/api/species", apiSpeciesHandler)

	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

	log.Println("Running on :8080")
	http.ListenAndServe(":8080", nil)
}
