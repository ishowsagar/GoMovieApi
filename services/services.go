package services

import (
	"context"
	"time"
)

// @ interface ( where all methods belongs to type => movie exists here and not exposed directly)
type MovieMethodStore interface {
	GetAllMovies()([]*Movie,error)
	CreateMovie(movie Movie) (*Movie,error)
}

// @ type for movie api data struct
// $ json tags to map fields with Json data format for api via reflection package
type Movie struct {
	ID string `json:"id"`
	Name string `json:"name"`
	Genre string `json:"genre"`
	Description string `json:"description"`
	Ratings float32 `json:"ratings"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// @ Methods belongs to (Type Movie) --> accessible through the Type
func (m Movie) GetAllMovies()([]*Movie,error) {
	// !querying requests with context for more flexible request with passed context
	ctx,cancel := context.WithTimeout(context.Background(),dbContextTimeOutDuration)
	defer cancel()

	// * Accessing db from var Db which stores Db connection returned from supply Func
	query := `
		select id,name,genre,description,ratings,created_at,updated_at from movies
	`
	MoviesRows,err := db.QueryContext(ctx,query)
	if err != nil {
		return nil,err
	}

	// var to store rows
	var movies []*Movie
	for MoviesRows.Next(){
		// access to each row/entry in table of data got from db query call
		var movie Movie
		// accessig each field and populating into movie var just created which matches movie queried data 
		err := MoviesRows.Scan(
			&movie.ID,
			&movie.Name,
			&movie.Genre,
			&movie.Description,
			&movie.Ratings,
			&movie.CreatedAt,
			&movie.UpdatedAt,
		)
		
		if err != nil {
			return nil,err
		}

	movies = append(movies, &movie)
	}
	return movies,err
}

// Create movie method
func (m Movie) CreateMovie(movie Movie) (*Movie,error) {
	ctx,cancel := context.WithTimeout(context.Background(),dbContextTimeOutDuration)
	defer cancel()

	query := `
		Insert into movies(name,genre,description,ratings,created_at,updated_at)
		Values($1,$2,$3,$4,$5,$6)
		returning *
	`

	_,err := db.ExecContext(ctx,query,
		movie.Name,
		movie.Genre,
		movie.Description,
		movie.Ratings,
		time.Now(),
		time.Now(),
	)
	if err !=nil {
		// using pointer as this data will be posted to db so need to be modified
		return nil,err
	}
	// if successfull query within context, return what create in the db
	return &movie,nil
}