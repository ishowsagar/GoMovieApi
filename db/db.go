// Provides connection to the Database to let Api call db and query things
package db

import (
	"database/sql"
	"fmt"
	"time"

	_ "github.com/jackc/pgconn"
	_ "github.com/jackc/pgx/stdlib"
	_ "github.com/jackc/pgx/v4"
	_ "github.com/lib/pq"
)

// type to store db connection
type Db struct {
	Db *sql.DB
}

//! empty instance --> stroes db Connection
var DbConnection = &Db{} //& defrenced cause --> pass & modify by pointer addr not copy

func ConnectToPostgresDB(connectionString string) (*Db,error) {

	// open connection - returns db & err
	db,err := sql.Open("pgx",connectionString)
	if err !=nil {
		fmt.Printf("failed to open database connection :%s",err)
	}

	// * configuring Db connection we have got from above func
	db.SetMaxIdleConns(7)
	db.SetMaxOpenConns(5)
	db.SetConnMaxLifetime(5 * time.Minute)

	// feeding in connection to instance we created above
	DbConnection.Db = db

	//# executing testing & debugging
	err = testDB(db)
	if err !=nil {
		return nil,err
	}	

	//! satisfying return types --> returning Instance of type struct Db --> storing connection
	return DbConnection,nil
	
}

//  func to test db connection
func testDB(db *sql.DB) error {
	err := db.Ping()
	if err !=nil {
		// when defining utils func --> return err directly --> so when invoked to store err and use it to wrap accordingly
		return err
	}
	fmt.Println("Database pinged successfully✅✅.")
	return nil
}