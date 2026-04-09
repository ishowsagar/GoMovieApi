package store

import (
	"context"
	"database/sql"
	"time"

	"github.com/ishowsagar/Go/movieApi/tokens"
)

// * types declaration
type DbTokenStore struct {
	Db *sql.DB
}

// ! all these function which creates db Connection just need to have supplied db connection
func NewDbTokenStore(db *sql.DB) *DbTokenStore {
	return &DbTokenStore{
		Db:db ,
	}
}

// @ TokenStore interface -> Stores all the methods that belongs to type --> DbTokenStore
type TokenStore interface {
	CreateToken(userID int,ttl time.Duration,scope string) (*tokens.Token,error)
	InsertToken (tokenPayload *tokens.Token) error
	DeleteAllTokenForUser(userID int, scope string) error
}


func (t *DbTokenStore) CreateToken(userID int,ttl time.Duration,scope string) (*tokens.Token,error) {

	// from tokens pckg generate token --> generates token struct data with applied hash 256 ✅✅
	token,err := tokens.GenerateToken(userID,ttl,scope)
	if err != nil {
		return nil,err
	}

	// using insert fnc to insert hashed token into token db --> does the db query call to insert token data into db ✅✅
	err = t.InsertToken(token)
	if err != nil {
		return nil,err
	}
	// return token ✅✅
	return token,err
}

func (t *DbTokenStore) InsertToken (tokenPayload *tokens.Token) error {

	// query to insert data into token
	ctx,cancel := context.WithTimeout(context.Background(),DbTimeout)
	defer cancel()

	// todo -> need to create token table from migrations sqlx
	// fixed - added token table
	query := `
		Insert into 
			tokens(hash,user_id,scope,expiry)
		Values 
			($1,$2,$3,$4)

	`

	// execute the db query by passing the recieved tokenPayload data
	_,err := t.Db.ExecContext(ctx,query,tokenPayload.Hash,tokenPayload.UserID,tokenPayload.Scope,tokenPayload.Expiry)
	if err != nil {
		return err
	}

	
	// if everything goes right and insertion is done
	return nil
}


func ( t *DbTokenStore) DeleteAllTokenForUser(userID int, scope string) error{
	ctx,cancel := context.WithTimeout(context.Background(),DbTimeout)
	defer cancel()

	query :=`
		Delete from 
			tokens
		where 
			user_id=$1 and scope=$2
	` 

	_,err := t.Db.ExecContext(ctx,query,userID,scope)
	if err != nil {
		return err
	}

	// if everything goes right and insertion is done
	return nil

}














