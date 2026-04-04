// Router for handling server redirected client requests
package router

import (
	"net/http"

	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
	"github.com/ishowsagar/Go/movieApi/controllers"
) 

func ServeRoutes() http.Handler {
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


	// @Main routes
	router.Route("/api",func(r chi.Router){
		// all routes defined below has parent route path "/api". For ex --> api/whateverRoutePath defined in below routes
		r.Get("/movies/all",controllers.GetAllMovies)
		r.Post("/movie/create",controllers.CreateMovie)
	})

	//! Returning router cause this satisfies --> http.Handler interface{serve}
	return router
}