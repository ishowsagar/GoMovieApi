// wires up api calls method to controller func which would be invoked by the client via routes
package controllers

import (
	"encoding/json"
	"fmt"
	"net/http"

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
	}

	// pass to method that creates movie ~ from decoded data into the movie struct
	// createdMovie,err := MovieStore.Movie.CreateMovie(movie) * before
	createdMovie,err := m.Store.CreateMovie(movie)

	// err handeling
	if err !=nil {
		utils.WriteJson(w,http.StatusBadRequest,utils.Envelop{"status":"unexpected error occurred at database end!"})
	}
	// send response back to the client
	utils.WriteJson(w,http.StatusOK,utils.Envelop{"data":createdMovie})
}