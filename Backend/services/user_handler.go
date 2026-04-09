package services

import (
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"
	"regexp"

	utils "github.com/ishowsagar/Go/movieApi/jsonUtils"
	"github.com/ishowsagar/Go/movieApi/store"
)

//! types declaration
//! registerUserRequest --> incoming JSON payload for user registration
type registerUserRequest struct {
	Username string `json:"username"` //* unique username for login
	Password string `json:"password"` //* plaintext password (will be hashed)
	Email    string `json:"email"` //* user email address
	Bio      string `json:"bio"` //* optional user bio
}

// @ interace where all methods belonsgs to type -> UserHandler that implements it
type UserHandlerInterface interface {
	ValidateUserRegisterRequest (regUser *registerUserRequest) error
	HandleRegisterUser(w http.ResponseWriter, req *http.Request)

}

type UserHandler struct {
	userStore store.UserStore //* database operations for users
	Logger *slog.Logger
}

//& type that stores interface of type UserHandler which has methods belongs to it 
	type UserHandlerInterfaceStore struct {
		UserHandlerIface UserHandlerInterface
	}


//! NewUserHandler --> constructor that creates user handler instance
func NewUserHandler(userStore store.UserStore,logger *slog.Logger) *UserHandler {
	return &UserHandler{
		userStore: userStore,
		Logger: logger,
	}
}

//! validateUserRegisterRequest --> server-side validation before saving to database
//? prevents invalid data from entering the system
func (h *UserHandler) ValidateUserRegisterRequest (regUser *registerUserRequest) error {
	//* checking if username is provided
	if regUser.Username == "" {
		h.Logger.Error("unexpected error occurred","error","must have a username")
		return errors.New("username is required")
	}
	
	if len(regUser.Username) > 50 {
		h.Logger.Error("unexpected error occurred","error","username length surpassing 50 characters")
		return errors.New("username cannot be greater than 50 characters")
	}

	if regUser.Email == "" {
		h.Logger.Error("unexpected error occurred","error","must have a email")
		return errors.New("email is required")
	}
	
	if regUser.Password == "" {
		h.Logger.Error("unexpected error occurred","error","password cannot be empty")
		return errors.New("password is required")
	}

	//! regex validation for email format
	emailRegexPattern := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	if !emailRegexPattern.MatchString(regUser.Email)  {
		//* test if email matches standard email pattern
		h.Logger.Error("unexpected error occurred","error","wrong email format,check your email address")
		return errors.New("Invalid email Xformat")
	} 

	return nil

}

//! HandleRegisterUser --> POST /users endpoint for creating new user accounts
func (h *UserHandler) HandleRegisterUser(w http.ResponseWriter, req *http.Request) {
	var r registerUserRequest //* holds incoming JSON data

	//* decode JSON body into struct
	err := json.NewDecoder(req.Body).Decode(&r)
	if err != nil {
		utils.WriteJson(w, http.StatusBadRequest, utils.Envelop{"error": "invalid request payload"})
		return
	}

	err = h.ValidateUserRegisterRequest(&r)
	if err != nil {
		utils.WriteJson(w, http.StatusBadRequest, utils.Envelop{"error": err.Error()})
		return
	}

	//* create User struct with validated data
	user := &store.User{
		Username: r.Username,
		Email:    r.Email,
	}
	if r.Bio != "" {
		user.Bio = r.Bio
	}

	//! hash the password using bcrypt (cost factor 12) - NEVER store plaintext passwords
	err = user.PasswordHash.SetUser(r.Password)
	if err != nil {
		utils.WriteJson(w, http.StatusInternalServerError, utils.Envelop{"error": "internal server error"})
		return
	}

	//* save user to database
	err = h.userStore.CreateUser(user)
	if err != nil {
		utils.WriteJson(w, http.StatusInternalServerError, utils.Envelop{"error": "internal server error"})
		return
	}

	//* 201 Created response with user data (password hash is excluded via json:"-" tag)
	utils.WriteJson(w, http.StatusCreated, utils.Envelop{"user": user})

}