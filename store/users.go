package store

import (
	"context"
	"crypto/sha256"
	"database/sql"
	"errors"
	"time"

	"golang.org/x/crypto/bcrypt"
)

//@ types declration
type Password struct {
	PlainText *string
	Hash []byte
}


// funcs that belongs to this type --> Password struct
func (p *Password) SetUser(plainTextPasswordString string) error {

	// need to return instance of &Password by setting hash and PlainText recieved from payload
	//# step 1 -> generate hash pass in form of []byte from bcrypt pckg
	hashedPassByte,err := bcrypt.GenerateFromPassword([]byte(plainTextPasswordString),12) // func demands pass in []byte() --> convert into byte and need cost factor
	if err != nil {
		return err
	}
	
	//# Step 2 -> set hash & and plain pass in the instance
	p.PlainText = &plainTextPasswordString
	p.Hash = hashedPassByte
	return nil
}


// func that comapares payload with hashed pass with bcrypt to check if user exists and pass matches
func (p *Password) ComparePassAndAuthenticate(plainTextPasswordString string) (bool,error) {
	err := bcrypt.CompareHashAndPassword([]byte(p.Hash),[]byte(plainTextPasswordString))
	if err != nil {
		switch {
		case errors.Is(err,bcrypt.ErrMismatchedHashAndPassword) :
			return false,nil
			default :
			return false,err
		}
	}
	return true,nil
}


// @ types delcaration for user
type User struct {
	Id int `json:"id"`
	Username string `json:"username"`
	Email string `json:"email"`
	PasswordHash Password `json:"-"`
	Bio string `json:"bio"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

const DbTimeout time.Duration = time.Second * 5

// @User Interface --> Stores all the methods that belongs to user Type --> type implements it
type UserStore interface {
	CreateUser(userPayload *User) (error)
	GetUserByUsername(username string) (*User,error)
	UpdateUser(user *User) (error)
	GetUserToken(scope,plainTextPassword string) (*User,error)
}


type DbUserStore struct {
	Db *sql.DB
}

func NewDbUserStore(db  *sql.DB) *DbUserStore {
	return &DbUserStore{
		Db: db,
	}
}

var anonyUser = &User{} // checks if incoming client is anon user or auth user
func ( u *User) CheckUser() bool {
	return u == anonyUser
}

// create user
func (d *DbUserStore) CreateUser(userPayload *User) (error) {

	ctx,cancel := context.WithTimeout(context.Background(),DbTimeout)
	defer cancel()

	// todo --> create migration to make a "User" table
	query := `
		Insert into Users
			(username,email,password_hash,bio)
		Values
			($1,$2,$3,$4)
		returning id,created_at,updated_at
	`

	res := d.Db.QueryRowContext(ctx,query,userPayload.Username,userPayload.Email,userPayload.PasswordHash.Hash,userPayload.Bio)

	// pupulate fields into pointer user payload giving back to userPayload
	err := res.Scan(
		&userPayload.Id,
		&userPayload.CreatedAt,
		&userPayload.UpdatedAt,
	 )
	if err != nil {
		return err
	}
	return nil
}


func (d *DbUserStore) GetUserByUsername(username string) (*User,error) {
	ctx,cancel := context.WithTimeout(context.Background(),DbTimeout)
	defer cancel()

	//  intiating instance of user data struct 
	user := &User{
		PasswordHash: Password{} ,
	}
	query :=`
		Select 
			id,username,email,password_hash,created_at,updated_at
		from
			users
		where
			username=$1
	`

	var passwordHash []byte

	err := d.Db.QueryRowContext(ctx,query,username).Scan(
		//! populating the resulting row into the user instance we have created 
		&user.Id,
		&user.Username,
		&user.Email,
		&passwordHash,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err != nil {
		return nil,err
	}
	user.PasswordHash.Hash = passwordHash
	return user,nil
}


func (d *DbUserStore) UpdateUser(user *User) error {
	ctx,cancel := context.WithTimeout(context.Background(),DbTimeout)
	defer cancel()	

	query := `
		update 
			users
		set
			username=$1,email=$2,bio=$3,updated_at=$4
		where
			id=$5
	`

	res,err := d.Db.ExecContext(ctx,query,user.Username,user.Email,user.Bio,time.Now(),user.Id)
	if err != nil {
		return err
	}

	rowsAffected,err := res.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return sql.ErrNoRows
	}
	return nil
}

// get token
func (d *DbUserStore) GetUserToken(scope,plainTextPassword string) (*User,error) {
	hashedPass := sha256.Sum256([]byte(plainTextPassword))
	ctx,cancel := context.WithTimeout(context.Background(),DbTimeout)
	defer cancel()	
	query := `
		Select 
			u.id,u.username,u.email,u.password_hash,u.bio,u.created_at,u.updated_at
		from
			users u
		Inner Join tokens t
		On
		t.user_id = u.id
		Where 
		t.hash=$1 And t.scope=$2 And t.expiry > $3
	`

	// intialioze user instance
	user := &User{
		PasswordHash: Password{},
	}
	// db query 
	res := d.Db.QueryRowContext(ctx,query,hashedPass[:],scope,time.Now())
	var passwordHash []byte

	// err handeling
	err := res.Scan(
		// resulting row scanning to fetch fields entry values by order and populate into
		&user.Id,
		&user.Username,
		&user.Email,
		&passwordHash,
		&user.Bio,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
	if err != nil {
		return nil,err
	}
	user.PasswordHash.Hash = passwordHash
	// return user
	return user,nil
}








