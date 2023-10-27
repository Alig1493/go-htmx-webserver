package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"io"
	"log"
	"net/http"
	"time"
)

type IMDBFilm struct {
	Rank        int      `json:"rank"`
	ID          string   `json:"id"`
	Title       string   `json:"title"`
	Description string   `json:"description"`
	IMDBD       string   `json:"imdbid"`
	Rating      string   `json:"rating"`
	Year        string   `json:"year"`
	Director    []string `json:"Director"`
	Writers     []string `json:"Writers"`
	Stars       []string `json:"Stars"`
	Genre       []string `json:"genre"`
	Images      []string `json:"images"`
}

type Film struct {
	Films *[]IMDBFilm
}

var myClient = &http.Client{Timeout: 20 * time.Second}

func fetch_movie_information() *[]IMDBFilm {
	url := "https://imdb-top-100-movies1.p.rapidapi.com/"

	req, _ := http.NewRequest("GET", url, nil)

	req.Header.Add("X-RapidAPI-Key", "7fe04a1c89msh80f96550ad521f0p112246jsnd37295424166")
	req.Header.Add("X-RapidAPI-Host", "imdb-top-100-movies1.p.rapidapi.com")
	req.Header.Add("Content-type", "application/json; charset=utf-8")

	res, error := myClient.Do(req)

	if error != nil {
		fmt.Println(error)
	}
	body, err := io.ReadAll(res.Body)
	if err != nil {
		log.Fatalf("Failed to read response body: %v", err)
	}

	defer res.Body.Close()

	var result []IMDBFilm
	json.Unmarshal(body, &result)

	return &result

}

func main() {
	// Hello world, the web server
	films := Film{
		Films: fetch_movie_information(),
	}

	indexHandler := func(w http.ResponseWriter, req *http.Request) {
		template := template.Must(template.ParseFiles("index.html"))
		error := template.Execute(w, films)
		if error != nil {
			log.Fatalln(error)
		}
	}

	add_film_handler := func(w http.ResponseWriter, req *http.Request) {
		title := req.PostFormValue("title")
		director := req.PostFormValue("director")
		log.Print(title)
		log.Print(director)
		template := template.Must(template.ParseFiles("index.html"))
		error := template.ExecuteTemplate(w, "film-list-element", IMDBFilm{Title: title, Director: []string{director}})
		if error != nil {
			log.Fatalln(error)
		}
	}

	http.HandleFunc("/", indexHandler)
	http.HandleFunc("/add-film/", add_film_handler)
	log.Fatal(http.ListenAndServe(":8000", nil))
}
