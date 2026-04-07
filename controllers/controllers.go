// wires up api calls method to controller func which would be invoked by the client via routes
package controllers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

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
	if idParam == "" {
		utils.WriteJson(w,http.StatusBadRequest,utils.Envelop{"status":"ID not found"})
		return
	}

	// make db call from interface passind id to invoke query
	retrievedMovie,err := m.Store.GetMovieByID(idParam)
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
	if idParam == "" {
		utils.WriteJson(w,http.StatusBadRequest,utils.Envelop{"status":"ID not found"})
		return
	}
	err := m.Store.DeleteMovieByID(idParam)
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