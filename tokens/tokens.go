package tokens

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/base32"
	"time"
)

// type declaration for the token data struct
type Token struct {
	PlainText string `json:"token"`
	Hash []byte `json:"-"`
	UserID int `json:"-"`
	Expiry time.Time `json:"expiry"`
	Scope string `json:"scope"`
}

const ScopeAuthForTokensScope = "authenticated"
// ! func that generates token with data
func GenerateToken(userID int,ttl time.Duration,scope string) (*Token,error) {
	// token instance with passed arg
	createdToken := &Token{
		UserID: userID,
		Expiry: time.Now().Add(ttl), //* current time + add ttl provided
		Scope: scope,
	}


	// generate random [] byte and read with "rand" pckg
	emptybyte := make([]byte,32)
	_,err := rand.Read(emptybyte) // converts to random byte
	if err != nil {
		return nil,err
	}

	// encoding plainText string to base32 for salt hashing
	createdToken.PlainText = base32.StdEncoding.WithPadding(base32.NoPadding).EncodeToString(emptybyte)
	// add sh56 algo on it for more salt
	hash := sha256.Sum256([]byte(createdToken.PlainText))
	// set token hash to be array not to be slice with conversion method
	createdToken.Hash = hash[:]
	// return token
	return createdToken,nil
}