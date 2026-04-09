package services

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"time"

	utils "github.com/ishowsagar/Go/movieApi/jsonUtils"
	"github.com/ishowsagar/Go/movieApi/store"
	"github.com/ishowsagar/Go/movieApi/tokens"
)

type TokenHandler struct {
	UserStore store.UserStore
	TokenStore store.TokenStore
	Logger *slog.Logger
}
//! createTokenRequest --> login credentials from client
type createTokenRequest struct {
	Username string `json:"username"` //* user's login name
	Password string `json:"password"` //* plaintext password to verify
}

// func that creates instance of type TokenHandler --> which will have methods belongs to it
func NewTokenHandler(UserStore store.UserStore ,TokenStore store.TokenStore,Logger *slog.Logger) *TokenHandler {
	return &TokenHandler{
		UserStore: UserStore,
		TokenStore: TokenStore,
		Logger: Logger,
	}
}

func (h *TokenHandler) HandleCreateToken(w http.ResponseWriter,r *http.Request) {

	//  check for tokenRequestingUser and decode body into it
	var requestingUser createTokenRequest
	err := json.NewDecoder(r.Body).Decode(&requestingUser)
	if err != nil {
		// ! this should be nil while feeding err
		h.Logger.Error("unexpected error occurred","error",err)
		utils.WriteJson(w,http.StatusBadRequest,utils.Envelop{"error":"failed to create token with passed data"})
		return 
	}

	//  server side validation for recieved and decoded payload
	switch {
	case requestingUser.Username == "" :
		h.Logger.Error("unexpected error occurred","error",err)
		utils.WriteJson(w,http.StatusBadRequest,utils.Envelop{"error":"username is mandatory"})
		return
	case requestingUser.Password == "" :
		h.Logger.Error("unexpected error occurred","error",err)
		utils.WriteJson(w,http.StatusBadRequest,utils.Envelop{"error":"password is mandatory"})
		return 
	}
	
	//  check if there is data existed for it when called getUserByUsername meth --> returns full user data struct
	retrievedUser,err := h.UserStore.GetUserByUsername(requestingUser.Username)
	if err != nil || retrievedUser == nil {		
		h.Logger.Error("unexpected error occurred","error",err)
		utils.WriteJson(w,http.StatusUnauthorized,utils.Envelop{"error":"invalid credentials"})
		return
	}

	// check if pass hash matches
	ok,err := retrievedUser.PasswordHash.ComparePassAndAuthenticate(requestingUser.Password)
	if err != nil {
		h.Logger.Error("unexpected error occurred","error",err)
		utils.WriteJson(w, http.StatusInternalServerError, utils.Envelop{"error": "internal server error"})
		return
	}

	if !ok {
		//? wrong password
		h.Logger.Error("unexpected error occurred","error",err)
		utils.WriteJson(w, http.StatusUnauthorized, utils.Envelop{"error": "invalid credentials"})
		return
	}
	// then create token with method that belongs to tokenStore
	token,err := h.TokenStore.CreateToken(retrievedUser.Id,15 * time.Minute,tokens.ScopeAuthForTokensScope)
	if err != nil {
		h.Logger.Error("unexpected error occurred","error",err)
		utils.WriteJson(w, http.StatusInternalServerError, utils.Envelop{"error": "internal server error"})
		return
	}

	// send res
	utils.WriteJson(w, http.StatusCreated, utils.Envelop{"auth_token": token })

}

