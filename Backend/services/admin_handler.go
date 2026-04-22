package services

import (
	"database/sql"
	"encoding/json"
	"errors"
	"net/http"
	"time"

	utils "github.com/ishowsagar/Go/movieApi/jsonUtils"
	"github.com/ishowsagar/Go/movieApi/store"
)

// method that serve admin
func(t *TokenHandler) HandleAdminSignup(w http.ResponseWriter,r *http.Request) {

	// incoming body req payload
	var adminSignupRequestPayload struct {
		Username string `json:"username"`
		Email string `json:"email"`
		Password string `json:"password"`
		Role string `json:"auth_role"`
	}

	// decode incoming body
	err :=json.NewDecoder(r.Body).Decode(&adminSignupRequestPayload)
	if err != nil {
		utils.WriteJson(w,http.StatusBadRequest,utils.Envelop{"error":"invalid admin signup request"}) 
		return
	}

	// atp adminReq holds data but needs validation
	switch {
	case adminSignupRequestPayload.Username == "" :
		utils.WriteJson(w,http.StatusBadRequest,utils.Envelop{"error":"username field is mandatory","json_expected":`{
			"username" : string,
			"password" : string,
			"email"    : string,
		}`})
		return
	case adminSignupRequestPayload.Password == "" :
		utils.WriteJson(w,http.StatusBadRequest,utils.Envelop{"error":"password field is mandatory","json_expected":`{
			"username" : string,
			"password" : string,
			"email"    : string,
		}`})
		return
	case adminSignupRequestPayload.Email == "" :
		utils.WriteJson(w,http.StatusBadRequest,utils.Envelop{"error":"email field is mandatory","json_expected":`{
			"username" : string,
			"password" : string,
			"email"    : string,
		}`})
		return
	case adminSignupRequestPayload.Role =="" :
		utils.WriteJson(w,http.StatusBadRequest,utils.Envelop{"error":"role field is mandatory","json_expected":`{
			"username" : string,
			"password" : string,
			"email"    : string,
		}`})
		return
	case adminSignupRequestPayload.Role != "admin" :
		utils.WriteJson(w,http.StatusBadRequest,utils.Envelop{"error":"role must be admin"})
		return		
	}

	// * now requested body is complete n validated

	existingUser,err := t.UserStore.GetUserByUsername(adminSignupRequestPayload.Username)
	if err!= nil && !errors.Is(err,sql.ErrNoRows){
		t.Logger.Error("failed to check existing admin user","error",err)
		utils.WriteJson(w,http.StatusInternalServerError,utils.Envelop{"error":"failed to check if user is already signed up"})
		return
	}

	// check if there is existing user in the database then could only create user
	if existingUser != nil {
		utils.WriteJson(w,http.StatusBadRequest,utils.Envelop{"error":"user already exists"})
		return
	}

	// generate user
	user := &store.User{
		Email: adminSignupRequestPayload.Email,
		Username: adminSignupRequestPayload.Username,
		// ! we need to store user password hash not concrete password
	}

	user.Bio = "admin"

	err  = user.PasswordHash.SetUser(adminSignupRequestPayload.Password) //* stores hash on instance where called
	if err != nil {
		t.Logger.Error("failed to hash admin passowrd -%v","err",err)
		utils.WriteJson(w,http.StatusInternalServerError,utils.Envelop{"error":"failed to hash admin password"})
		return
	}

	err  = t.UserStore.CreateUser(user)
	if err != nil {
		t.Logger.Error("failed to create admin user", "error", err)
		utils.WriteJson(w,http.StatusInternalServerError,utils.Envelop{"error":"failed to create admin role user"})
		return
	}

	// now attached context on it
	token,err := t.TokenStore.CreateToken(user.Id,15 * time.Minute,"admin_authenticated")
	if err != nil {
		t.Logger.Error("failed to create admin auth token", "error", err)
		utils.WriteJson(w,http.StatusInternalServerError,utils.Envelop{"error":"failed to create token for admin."})
		return
	}
	
	utils.WriteJson(w,http.StatusAccepted,utils.Envelop{"status":"admin created successfully✅✅","admin":token,})
}


// serve login
func(t *TokenHandler) HandleAdminLogin(w http.ResponseWriter,r *http.Request) {

	// incoming body req payload
	var adminLoginRequestPayload struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}

	// decode incoming body
	err :=json.NewDecoder(r.Body).Decode(&adminLoginRequestPayload)
	if err != nil {
		utils.WriteJson(w,http.StatusBadRequest,utils.Envelop{"error":"invalid admin login request"}) 
		return
	}

	// atp adminLoginReq holds data but needs validation
	
	if adminLoginRequestPayload.Username == "" {
		utils.WriteJson(w,http.StatusBadRequest,utils.Envelop{"error":"username field is mandatory","json_expected":`{
			"username" : string,
			"password" : string,
			}`})
		return
	}
	if adminLoginRequestPayload.Password == "" {
		utils.WriteJson(w,http.StatusBadRequest,utils.Envelop{"error":"password field is mandatory","json_expected":`{
			"username" : string,
			"password" : string,
			}`})
		return		
	}
	

	// * now requested body is complete n validated

	existingUser,err := t.UserStore.GetUserByUsername(adminLoginRequestPayload.Username)
	if err!= nil && !errors.Is(err,sql.ErrNoRows){
		t.Logger.Error("failed to check existing admin user","error",err)
		utils.WriteJson(w,http.StatusInternalServerError,utils.Envelop{"error":"failed to check if user is already signed up"})
		return
	}
	if errors.Is(err, sql.ErrNoRows) || existingUser == nil {
		utils.WriteJson(w,http.StatusUnauthorized,utils.Envelop{"error":"invalid credentials"})
		return
	}

	// check pass hash
	ok,err := existingUser.PasswordHash.ComparePassAndAuthenticate(adminLoginRequestPayload.Password) //* stores hash on instance where called

	if err != nil {
		t.Logger.Error("failed to check admin password", "error", err)
		utils.WriteJson(w,http.StatusInternalServerError,utils.Envelop{"error":"failed to check admin password"})
		return
	}
	if !ok {
		utils.WriteJson(w,http.StatusUnauthorized,utils.Envelop{"error":"invalid credentials"})
		return
	}

	// // generating jwt
	// // # step 1 --> generate jwt token with jwt.NewMapWithClaims method passing singin method & claims data struct
	// token := jwt.NewWithClaims(jwt.SigningMethodHS256,jwt.MapClaims{
	// 	// adding map data which need to pass to the token
	// 	"userId":existingUser.Id,
	// 	"expiry": time.Now().Add(15 * time.Minute).Unix(),
	// })

	// tokenString,err := token.SignedString([]byte("secret"))
	// if err!= nil {
	// 	t.Logger.Error("failed to issue token for req")
	// 	return
	// }
	token,err := t.TokenStore.CreateToken(existingUser.Id,15 * time.Minute,"admin_authenticated") // inserts created token into the database
	if err != nil {
		t.Logger.Error("failed to create admin login token", "error", err)
		utils.WriteJson(w,http.StatusInternalServerError,utils.Envelop{"error":"failed to create token for admin"})
		return
	}

	utils.WriteJson(w,http.StatusOK,utils.Envelop{"status":"admin login successful","admin-token":token})

}

