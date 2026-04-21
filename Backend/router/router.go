// Router for handling server redirected client requests
package router

import (
	"net/http"

	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
	"github.com/ishowsagar/Go/movieApi/controllers"
	mw "github.com/ishowsagar/Go/movieApi/middleware"
	"github.com/ishowsagar/Go/movieApi/services"
)

// @ this function serves route and has access to Application which stores db Conn
func ServeRoutes(th *services.TokenHandler,mw mw.UserMiddleware,usrHandlerIfaceStore services.UserHandlerInterfaceStore ) http.Handler {
	router := chi.NewRouter()

	// * Configuring router to have mw & enabled access to domains N meths
	router.Use(middleware.Recoverer)
	router.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"http://localhost:3000", "http://127.0.0.1:3000", "http://localhost:5173", "http://127.0.0.1:5173"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300,
	}))

	// * Health check route to check Standard Go App is running or not
	router.Get("/health",func(w http.ResponseWriter,r *http.Request) {
		
		responseMsg := "Go app is running fine🔋🔋..."

		w.Header().Set("Content-type","application/json")
		w.WriteHeader(http.StatusOK)
		_,err := w.Write([]byte(responseMsg)) // _ to not to store returned intCode but err

		if err !=nil {
			w.WriteHeader(http.StatusBadRequest)
		}
	})

	// & Accessing controller methods by type that they belongs to Movie type
	var movie services.Movie
	controllers := controllers.NewMovieMethodStore(movie)

	// @Main routes
	router.Route("/api",func(r chi.Router){
		// all routes defined below has parent route path "/api". For ex --> api/whateverRoutePath defined in below routes
		// tokenStore :=  
		// userStore :=
		// tokenHandler := services.NewTokenHandler()

		// @ Chaining router with mw calls to let them use it
		r.Use(mw.Authenticate) // uncomment to use auth for protection enabled
		//! Protected Routes (uncomment to check when login/signup is implemented)
		r.Post("/users/register",usrHandlerIfaceStore.UserHandlerIface.HandleRegisterUser)
		r.Post("/tokens/authentication",th.HandleCreateToken)
		r.Get("/movies/all",mw.RequiresAuthorization(controllers.GetAllMovies))
		r.Get("/movies/movie/{id}", mw.RequiresAuthorization(controllers.GetMovieByID)) //* urlParam is read when ending slug is in format of /{slug}
		r.Post("/movie/create",mw.RequiresAuthorization(controllers.CreateMovie))
		r.Put("/movies/movie/update/{id}",mw.RequiresAuthorization(controllers.UpdateMovieByID))
		r.Delete("/movies/movie/delete/{id}",mw.RequiresAuthorization(controllers.DeleteMovieByID))
		r.Delete("/movies/all/delete",mw.RequiresAuthorization(controllers.DeleteAllMovies))
	
		// r.Post("/users/register",usrHandlerIfaceStore.UserHandlerIface.HandleRegisterUser)
		// r.Post("/tokens/authentication",th.HandleCreateToken)
		// r.Get("/movies/all",controllers.GetAllMovies)
		// r.Get("/movies/movie/{id}", controllers.GetMovieByID) //* urlParam is read when ending slug is in format of /{slug}
		// r.Get("/movies/movie/query/name",controllers.GetMovieByQueryParam) //* ? query param after ?
		// r.Get("/movies/movie/query/genre",controllers.GetMoviesByGenre)
		// r.Get("/movies/movie/query/ratings",controllers.GetMoviesByRatingsQP)
		// r.Get("/movies/{genre}/query",controllers.GetMoviesByLimit)
		// r.Post("/movie/create",controllers.CreateMovie)
		// r.Put("/movies/movie/update/{id}",controllers.UpdateMovieByID)
		// r.Delete("/movies/movie/delete/{id}",controllers.DeleteMovieByID)
		// r.Delete("/movies/all/delete",controllers.DeleteAllMovies)
		
	})

	//! Returning router cause this satisfies --> http.Handler interface{serve}
	return router
}
