package main

import (
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

// Movie model
// Represents a film with ID, ISBN, title, and director information
type Movie struct {
	ID       string    `json:"id"`
	Isbn     string    `json:"isbn"`
	Title    string    `json:"title"`
	Director *Director `json:"director"`
}

// Director model
// Represents a director with first and last name
type Director struct {
	Firstname string `json:"firstname"`
	Lastname  string `json:"lastname"`
}

// movies slice to simulate a database
var movies []Movie

// getMovies handles GET requests to retrieve all movies
func getMovies(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(movies)
}

// getMovie handles GET requests to retrieve a single movie by ID
func getMovie(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r) // capture route parameters
	for _, item := range movies {
		if item.ID == params["id"] {
			json.NewEncoder(w).Encode(item)
			return
		}
	}
	// Respond with error message if movie not found
	json.NewEncoder(w).Encode(map[string]string{"message": "Movie not found"})
}

// createMovie handles POST requests to add a new movie
func createMovie(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var movie Movie
	_ = json.NewDecoder(r.Body).Decode(&movie)     // decode request body into movie struct
	movie.ID = strconv.Itoa(rand.Intn(1_000_000))  // generate random ID
	movies = append(movies, movie)                 // add to slice
	json.NewEncoder(w).Encode(movie)               // return created movie
}

// updateMovie handles PUT requests to update an existing movie by ID
func updateMovie(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r) // capture route parameters
	for index, item := range movies {
		if item.ID == params["id"] {
			// remove the existing movie
			movies = append(movies[:index], movies[index+1:]...)

			var movie Movie
			if err := json.NewDecoder(r.Body).Decode(&movie); err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}
			movie.ID = params["id"]        // preserve original ID
			movies = append(movies, movie) // add updated movie
			json.NewEncoder(w).Encode(movie)
			return
		}
	}
	// Respond with error message if movie not found
	json.NewEncoder(w).Encode(map[string]string{"message": "Movie not found"})
}

// deleteMovie handles DELETE requests to remove a movie by ID
func deleteMovie(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r) // capture route parameters
	for index, item := range movies {
		if item.ID == params["id"] {
			movies = append(movies[:index], movies[index+1:]...) // remove movie
			break
		}
	}
	json.NewEncoder(w).Encode(movies) // return remaining movies
}

// main function to start the server and initialize routes
func main() {
	r := mux.NewRouter()

	// initial data
	movies = append(movies, Movie{ID: "1", Isbn: "438277", Title: "Movie One", Director: &Director{Firstname: "John", Lastname: "Doe"}})
	movies = append(movies, Movie{ID: "2", Isbn: "277438", Title: "Movie Two", Director: &Director{Firstname: "Smith", Lastname: "Bob"}})

	// define routes and handlers
	r.HandleFunc("/movies", getMovies).Methods("GET")
	r.HandleFunc("/movies/{id}", getMovie).Methods("GET")
	r.HandleFunc("/movies", createMovie).Methods("POST")
	r.HandleFunc("/movies/{id}", updateMovie).Methods("PUT")
	r.HandleFunc("/movies/{id}", deleteMovie).Methods("DELETE")

	fmt.Println("Starting server at port 8000")
	log.Fatal(http.ListenAndServe(":8000", r))
}
