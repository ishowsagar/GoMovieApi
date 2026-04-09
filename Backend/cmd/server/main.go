package main

import (
	"fmt"
	"log/slog"
	"os"
	"time"

	"net/http"

	"github.com/ishowsagar/Go/movieApi/db"
	mw "github.com/ishowsagar/Go/movieApi/middleware"
	"github.com/ishowsagar/Go/movieApi/router"
	"github.com/ishowsagar/Go/movieApi/services"
	"github.com/ishowsagar/Go/movieApi/store"
	"github.com/joho/godotenv"
)

// @Types declarations
// whatever needs in application method would be recieved through pointer instance
type Application struct {
	Logger *slog.Logger
	Config config
	MovieStore services.MovieStore
	TokenHandler *services.TokenHandler
	// UserHandler *services.UserHandler
	MiddleWare mw.UserMiddleware
	UserHandlerInterface services.UserHandlerInterfaceStore
}
type config struct {
	PORT string
}

// @ Imp utils inventory 

func(a *Application) IntializeServer() error {
	chiRouter := router.ServeRoutes(a.TokenHandler,a.MiddleWare,a.UserHandlerInterface)
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
		return
	}

	// & logger for flexible debugging
	logger := slog.New(slog.NewTextHandler(os.Stdout,&slog.HandlerOptions{
		Level: slog.LevelInfo,
	}))
	slog.SetDefault(logger) //* tuning up logger to be in action as in Globally accessible state

	// accessing env vars for use
	port := os.Getenv("PORT")
	DBConnStr := os.Getenv("DB_CONN_STR")
	
	
	databaseConnection,err := db.ConnectToPostgresDB(DBConnStr)
	if err !=nil {
		// fmt.Printf("fmtLog : failed to engine up postgres database - %s",err)
		slog.Warn("failed to engine up postgres database","err",err)
		return
	}

	// !remember => when you need to pass something,you could to everywhere via paramter lik we are passing Logger
	
	//@ enabling stores into action by providing'em db conn --> consumed by router
	UserStore := store.NewDbUserStore(databaseConnection.Db)
	TokenStore := store.NewDbTokenStore(databaseConnection.Db)
	TokenHandler := services.NewTokenHandler(UserStore,TokenStore,logger)
	UserHandler := services.NewUserHandler(UserStore,logger)
	UserHandlerInterface := services.UserHandlerInterfaceStore{UserHandlerIface: UserHandler}
	MiddleWareHandler := mw.UserMiddleware{UserStore: *UserStore,Logger: logger}

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
		TokenHandler: TokenHandler,
		MiddleWare: MiddleWareHandler,
		UserHandlerInterface: UserHandlerInterface,
		Logger: logger,
	}

	err = app.IntializeServer()
	if err !=nil {
		logger.Error("server error","err","failed to fuel up the server")
		fmt.Printf("failed to start server")
		return
	}
	fmt.Println("Go app has started🚀...")
}

 


