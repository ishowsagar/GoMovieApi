package main

import (
	"fmt"
	"os"
	"time"

	"net/http"

	"github.com/ishowsagar/Go/movieApi/db"
	"github.com/ishowsagar/Go/movieApi/router"
	"github.com/ishowsagar/Go/movieApi/services"
	"github.com/joho/godotenv"
)

// @Types declarations
type Application struct {
	Config config
	MovieStore services.MovieStore
}
type config struct {
	PORT string
}

// @ Imp utils inventory 

func(a *Application) IntializeServer() error {
	chiRouter := router.ServeRoutes()
	server := &http.Server{
		Addr: fmt.Sprintf(":%s",a.Config.PORT),
		ReadTimeout: 4 * time.Second,
		WriteTimeout: 3 * time.Second,
		IdleTimeout: 5 * time.Second,
		Handler: chiRouter,
	}

	return server.ListenAndServe()
	
	
}

// main func 
func main() {

	// loading .env file 
	err := godotenv.Load()
	if err !=nil {
		fmt.Printf("failed to load env file")
	}

	// accessing env vars for use
	port := os.Getenv("PORT")
	DBConnStr := os.Getenv("DB_CONN_STR")

	databaseConnection,err := db.ConnectToPostgresDB(DBConnStr)
	if err !=nil {
		fmt.Printf("failed to load db connection")
	}
	defer func(){
		if err == nil {
			databaseConnection.Db.Close() //deferred to be invoked at the end when all sorrouding func gets invoked --> cleanup conn at en
		}
	}() 
	

	// & this db stores db instance which has sqlConnection --> supplied to api func
	
	
	
	// ! creating Defrenced &instance from type Application type struct and to access methods defined on it
	app := &Application{
		Config: config{
			PORT:port,
		},
		MovieStore : services.SupplyDbConnectionToAPI(databaseConnection.Db),//& instansiates the model struct and this fnc also assigns passed dbConnection.db stored from db type returned from db func to the db var (now holds the actual db connection) used by api
	}

	err = app.IntializeServer()
	if err !=nil {
		fmt.Printf("failed to start server")
	}
	fmt.Println("Go app has started🚀...")
}





