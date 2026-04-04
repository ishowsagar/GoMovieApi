package services

import (
	"database/sql"
	"time"
)

// !db type that stores connection from db pkg func returned type
var db *sql.DB
const dbContextTimeOutDuration = 4 * time.Second

// stores movie type
type MovieStore struct {
	Movie Movie
}

// ! this func creates instance of store and takes in db connection --> when invoked --> passed in Db.db ( not concrete as we are directly accesing stored connection in the Db type struct)
func SupplyDbConnectionToAPI(DbConnection *sql.DB) MovieStore {
	// # assigns connection from Db.Db to var db which would be used by api to query db  
	db = DbConnection
	// this returns instance of MovieStore and also assigns dbConn to Api package's var db --> makes use of it to query db calls 
	return MovieStore{}
}