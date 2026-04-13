// wires up api calls method to controller func which would be invoked by the client via routes
package controllers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/go-chi/chi/v5"
	utils "github.com/ishowsagar/Go/movieApi/jsonUtils"
	"github.com/ishowsagar/Go/movieApi/services"
)

// @ accessing Store
var MovieStore services.MovieStore

// storing Interface 
type MovieMethodStore struct {
	Store services.MovieMethodStore
}
//  func that returns instance of type that holds interface passin in type that implements it
func NewMovieMethodStore(movieInterfaceImplemenationType services.MovieMethodStore) MovieMethodStore {
	return MovieMethodStore{
		Store: movieInterfaceImplemenationType,
	}
}

// ! All Controller~Handler methods belongs to type MovieStore interface which have all the apu calls methods
func (m MovieMethodStore) GetAllMovies(w http.ResponseWriter,r *http.Request)  {

	// RetrievedMovies,err := MovieStore.Movie.GetAllMovies() - * before
	RetrievedMovies,err := m.Store.GetAllMovies()
	if err != nil {
		fmt.Printf("failed to get all movies,Please try again later")
		w.WriteHeader(http.StatusBadRequest)
		utils.WriteJson(w,http.StatusBadRequest,utils.Envelop{"status":"failed to get all movies."})
		return
	}

	//  if client hit no error --> send resp✅✅
	utils.WriteJson(w,http.StatusOK,utils.Envelop{"data":RetrievedMovies})


}

func (m MovieMethodStore) CreateMovie(w http.ResponseWriter,r *http.Request) {

	// decode req.body
	var movie services.Movie

	err := json.NewDecoder(r.Body).Decode(&movie) //reflection does the population work under the hood to examine json struct tags & populate data 
	if err !=nil {
		utils.WriteJson(w,http.StatusBadRequest,utils.Envelop{"status":"unexpected error occurred, please send correct movie format data!"})
		return
	}
	switch {
	case movie.Name == "" :
		utils.WriteJson(w,http.StatusBadRequest,utils.Envelop{"status":"failed to update movie, name field is mandatory!"})
		return
	case movie.Genre == "" :
		utils.WriteJson(w,http.StatusBadRequest,utils.Envelop{"status":"failed to update movie, genre field is mandatory!"})
		return
	case movie.Description == "" :
		utils.WriteJson(w,http.StatusBadRequest,utils.Envelop{"status":"failed to update movie, description field is mandatory!"})
		return
	case movie.Ratings == 0.0 :
		utils.WriteJson(w,http.StatusBadRequest,utils.Envelop{"status":"failed to update movie, ratings field is mandatory!"})
		return
	}

	// pass to method that creates movie ~ from decoded data into the movie struct
	// createdMovie,err := MovieStore.Movie.CreateMovie(movie) * before
	createdMovie,err := m.Store.CreateMovie(movie)

	// err handeling
	if err !=nil {
		utils.WriteJson(w,http.StatusBadRequest,utils.Envelop{"status":"unexpected error occurred at database end!"})
		return
	}
	// send response back to the client
	utils.WriteJson(w,http.StatusOK,utils.Envelop{"data":createdMovie})
}

//! get movie by its id --> api/movies/movie/{id} /{slug being id}
func (m MovieMethodStore) GetMovieByID(w http.ResponseWriter,r *http.Request) {
	
	// retrieve id from url passed by client
	idParam := chi.URLParam(r,"id")
	cleanIDSlug := strings.TrimSpace(idParam) 
	if idParam == "" {
		utils.WriteJson(w,http.StatusBadRequest,utils.Envelop{"status":"ID not found"})
		return
	}

	// make db call from interface passind id to invoke query
	retrievedMovie,err := m.Store.GetMovieByID(cleanIDSlug)
	if err != nil {
		utils.WriteJson(w,http.StatusBadRequest,utils.Envelop{"status":"Movie not found, make sure ID is correct"})
		return
	}

	// send resp back to client ✅✅
	utils.WriteJson(w,http.StatusOK,utils.Envelop{"movie":retrievedMovie})
}

// ! Update movie with id and body to update with
func(m MovieMethodStore) UpdateMovieByID(w http.ResponseWriter,r *http.Request) {
	
	// decode incoming body from req from client
	var movie services.Movie
	// store in instance to pupulate that decoded data struct
	err := json.NewDecoder(r.Body).Decode(&movie)
	if err != nil {
		utils.WriteJson(w,http.StatusBadRequest,utils.Envelop{"status":"failed to update with provided json body, retry with correct data format"})
		return
	}
	
	// successfully populated recieved data from client into the var movie instance 
	
	// but validating if every field in movie var which is pupulated is not
	switch {
	case movie.Name == "" :
		utils.WriteJson(w,http.StatusBadRequest,utils.Envelop{"status":"failed to update movie, name field is mandatory!"})
		return
	case movie.Genre == "" :
		utils.WriteJson(w,http.StatusBadRequest,utils.Envelop{"status":"failed to update movie, genre field is mandatory!"})
		return
	case movie.Description == "" :
		utils.WriteJson(w,http.StatusBadRequest,utils.Envelop{"status":"failed to update movie, description field is mandatory!"})
		return
	case movie.Ratings == 0.0 :
		utils.WriteJson(w,http.StatusBadRequest,utils.Envelop{"status":"failed to update movie, ratings field is mandatory!"})
		return
	}
	
	// inject both id and body to update method
	idParam := chi.URLParam(r,"id")
	if idParam == "" {
		utils.WriteJson(w,http.StatusBadRequest,utils.Envelop{"status":"ID not found,Please provide ID of movie which was going to update"})
		return
	}
	// send res ✅✅
	utils.WriteJson(w,http.StatusAccepted,utils.Envelop{"status":"movie updated successfully."})
}


// delete movie by id --> /api/movies/movie/delete/{id}
func (m MovieMethodStore) DeleteMovieByID(w http.ResponseWriter,r *http.Request) {

	idParam := chi.URLParam(r,"id")
	cleanIDSlug := strings.TrimSpace(idParam)
	if idParam == "" {
		utils.WriteJson(w,http.StatusBadRequest,utils.Envelop{"status":"ID not found"})
		return
	}
	err := m.Store.DeleteMovieByID(cleanIDSlug)
	if err != nil {
		utils.WriteJson(w,http.StatusBadRequest,utils.Envelop{"status":"failed to delete movie, movie not found with passed ID"})
		return
	}

	utils.WriteJson(w,http.StatusAccepted,utils.Envelop{"status" : "successfully deleted movie."})
}

// delete all movies from the database
func (m MovieMethodStore) DeleteAllMovies(w http.ResponseWriter,r *http.Request) {

	err := m.Store.DeleteAllMovies()
	if err != nil {
		log.Fatalf("failed to delete all movies - %s",err)
		utils.WriteJson(w,http.StatusBadRequest,utils.Envelop{"status":"failed to delete all movies, unexpected error occured"})
		return
	}
	// query being successfull
	utils.WriteJson(w,http.StatusAccepted,utils.Envelop{"status": "Wiped all movies data from the database"})
}

// get movie by query param -> /api/movies/query?movie={}
func (m MovieMethodStore) GetMovieByQueryParam(w http.ResponseWriter,r *http.Request) {

	queryParam := r.URL.Query().Get("movie") // this is query ? param, not URlparam
	if queryParam == "" {
		utils.WriteJson(w,http.StatusBadRequest,utils.Envelop{"error":"movie query is missing"})
		return
	}	

	retrievedMovie,err := m.Store.GetMovieByQueryParams(queryParam)
	if err != nil {
		// log.Fatal("failed to get movie from db")
		er := err.Error()
		utils.WriteJson(w,http.StatusBadRequest,utils.Envelop{"error":er})
		return
	}
	utils.WriteJson(w,http.StatusOK,utils.Envelop{"movie":retrievedMovie})
}

// get movies of genre type qp --> /api/movies/query/genre={}
func (m MovieMethodStore) GetMoviesByGenre(w http.ResponseWriter,r *http.Request) {

	// retrieve genre from the url query
	genreQP := r.URL.Query().Get("genre") 
	if genreQP == "" {
		utils.WriteJson(w,http.StatusBadRequest,utils.Envelop{"error" : "genre query is missing!."})
		return
	}
	// make call to method to retrieve movies by making a db call
	retreivedMovies,err := m.Store.GetMoviesByGenreQP(genreQP)
	if err != nil {
		// err handeling
		utils.WriteJson(w,http.StatusBadRequest,utils.Envelop{"error" : "unknown genre, please pass correct genre!."})
		return
	}
	// send resp to client
	utils.WriteJson(w,http.StatusOK,utils.Envelop{"movies" : retreivedMovies})
}

// get movies of ratings qp --> /api/movies/query/ratings?={}
func (m MovieMethodStore) GetMoviesByRatingsQP(w http.ResponseWriter,r *http.Request) {

	// retrieve ratings from the url query
	ratingsQP := r.URL.Query().Get("rating") 
	if ratingsQP == "" {
		utils.WriteJson(w,http.StatusBadRequest,utils.Envelop{"error" : "rating query is missing!."})
		return
	}
	// ! we are recieving ratings as string so we need to convert it to uint
	parsedRating,err := strconv.Atoi(ratingsQP)
	if err != nil {
		// err handeling
		utils.WriteJson(w,http.StatusBadRequest,utils.Envelop{"error" : "please pass numerical rating value only."})
		return
	}

	// make call to method to retrieve movies by making a db call
	retreivedMovies,err := m.Store.GetMoviesByRatingsQP(uint(parsedRating))
	if err != nil {
		// err handeling
		utils.WriteJson(w,http.StatusNotFound,utils.Envelop{"error" : "no movies found with ratings check"})
		return
	}
	// send resp to client
	utils.WriteJson(w,http.StatusOK,utils.Envelop{"movies" : retreivedMovies})
}

// get limited no of movies by specifying genre and offset return limit -> /api/movies/{genre}?"limit"={}
func (m MovieMethodStore) GetMoviesByLimit(w http.ResponseWriter,r *http.Request) {

	rawGenreURlParamSlugWithSpacesandTrailingStuff := chi.URLParam(r,"genre") 
	genre := strings.TrimSpace(rawGenreURlParamSlugWithSpacesandTrailingStuff)
	if genre == "" {
		utils.WriteJson(w,http.StatusBadRequest,utils.Envelop{"error":"must provide genre slug, expected /api/movies/{genre}!."})
		return
	}

	stringyfiedLimitOffset := r.URL.Query().Get("limit") // stringed version in url query param
	parsedLimitOffset,err := strconv.Atoi(stringyfiedLimitOffset) // converting it to numerical val if could be possible 
	if err != nil {
		utils.WriteJson(w,http.StatusBadRequest,utils.Envelop{"error":"please provide correct slug only and with correct format"})
		return
	}
	retrievedMovies,err := m.Store.GetMoviesBYLimit(genre,uint8(parsedLimitOffset)) //* genre,offset limit
	// first checking if it has returned any rows or not before general err handeling
	if err == sql.ErrNoRows {
		utils.WriteJson(w,http.StatusBadRequest,utils.Envelop{"error":err.Error()}) // have to convert err to string first to send to resp
		return
	}
	if err!= nil{
		utils.WriteJson(w,http.StatusInternalServerError,utils.Envelop{"error":"failed to fetch movies by genre slug!."})
		return
	} 

	utils.WriteJson(w,http.StatusOK,utils.Envelop{"movies":retrievedMovies})
}