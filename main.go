package main

import (
	"encoding/json"
	"fmt"
	"go-htmx-webserver/pagination"
	"html/template"
	"io"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/bxcodec/faker"
)

type IMDBFilm struct {
	Rank        int      `json:"rank" faker:"boundary_start=5, boundary_end=10"`
	ID          string   `json:"id" faker:"UUIDDigit"`
	Title       string   `json:"title" faker:"name"`
	Description string   `json:"description" faker:"sentence"`
	IMDBD       string   `json:"imdbid" faker:"word"`
	Rating      string   `json:"rating" faker:"amount"`
	Year        string   `json:"year" faker:"year"`
	Director    []string `json:"Director" faker:"slice_len=2"`
	Writers     []string `json:"Writers" faker:"slice_len=2"`
	Stars       []string `json:"Stars" faker:"slice_len=2"`
	Genre       []string `json:"genre" faker:"slice_len=2"`
	Images      []string `json:"images" faker:"url, slice_len=2"`
}

type Error struct {
	Message string `json:"message"`
}

type Data struct {
	Films *[]IMDBFilm
	Pages []pagination.Page
}

var myClient = &http.Client{Timeout: 20 * time.Second}

func fetch_movie_information() (int, *[]IMDBFilm) {
	url := "https://imdb-top-100-movies1.p.rapidapi.com/"

	req, _ := http.NewRequest("GET", url, nil)

	req.Header.Add("X-RapidAPI-Key", "7fe04a1c89msh80f96550ad521f0p112246jsnd37295424166")
	req.Header.Add("X-RapidAPI-Host", "imdb-top-100-movies1.p.rapidapi.com")
	req.Header.Add("Content-type", "application/json; charset=utf-8")

	res, error := myClient.Do(req)

	log.Println("The status code we got is:", res.StatusCode)
	log.Println("The status code text we got is:", http.StatusText(res.StatusCode))

	if error != nil {
		log.Println(error)
	}

	body, err := io.ReadAll(res.Body)
	if err != nil {
		log.Fatalf("Failed to read response body: %v", err)
	}

	defer res.Body.Close()
	var result []IMDBFilm

	if res.StatusCode > 299 || res.StatusCode < 200 {
		log.Println(http.StatusText(res.StatusCode))
		var error Error
		json.Unmarshal(body, &error)
		log.Printf("%+v", error)

		for i := 0; i < 4; i++ {
			single := IMDBFilm{}
			err := faker.FakeData(&single)
			if err != nil {
				fmt.Println("An error occured when generating fake data due to: ", err)
			}

			result = append(result, single)
		}

		log.Printf("%+v", result)
		return len(result), &result
	}

	json.Unmarshal(body, &result)
	log.Print("Body: ")
	log.Println(string(body[:]))
	log.Println(result)

	return len(result), &result

}

func main() {
	// Hello world, the web server
	// films := Film{
	// 	Films: fetch_movie_information(),
	// }
	movie_list_length, movies := fetch_movie_information()

	templates := []string{
		"index.html",
		"pagination.html",
		"film-list.html",
		"add-film.html",
	}

	indexHandler := func(w http.ResponseWriter, req *http.Request) {
		query := req.URL.Query()
		starting_page_string := query.Get("starting_page")
		if starting_page_string == "" {
			starting_page_string = "1"
		}

		starting_page, err := strconv.Atoi(starting_page_string)
		log.Println("Starting page: ", starting_page)

		if err != nil {
			log.Println("Unable to convert string to integer")
			starting_page = 1
		}
		template := template.Must(template.ParseFiles(templates...))
		start_index := (starting_page - 1) * pagination.PSIZE
		end_index := ((starting_page - 1) * pagination.PSIZE) + pagination.PSIZE
		log.Println("Starting page element:", start_index)
		log.Println("Ending Page element: ", end_index)

		if end_index > movie_list_length {
			end_index = movie_list_length
		}

		sliced_movies := (*movies)[start_index:end_index]
		log.Println("Sliced movies: ", sliced_movies)
		films := Data{
			Films: &sliced_movies,
			Pages: pagination.Pager(starting_page, movie_list_length),
		}
		log.Println(films)
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
