package mw

import (
	"context"
	"log/slog"
	"net/http"
	"strings"

	utils "github.com/ishowsagar/Go/movieApi/jsonUtils"
	"github.com/ishowsagar/Go/movieApi/store"
	"github.com/ishowsagar/Go/movieApi/tokens"
)

// types declaration
type UserMiddleware struct {
	UserStore store.DbUserStore
	Logger *slog.Logger
}

const userContextKey string= "user"

// func that set user in ctx via new req state with attached ctx
func SetUser(user *store.User,r *http.Request) *http.Request {
	// * creating context on source r.Context existing context and adding this key-val new attached ctx
	ctx := context.WithValue(r.Context(),userContextKey,user) //* new context state
	return r.WithContext(ctx) //* New req state w/ context attched to it
}  

// func that get user from req's attached ctx with ANyctx.Value(keyToCheckFor).(againstType)
func GetUser(r *http.Request) *store.User {
	// comma ok to check if something value is coming alright, ok is bool for it
	user,ok := r.Context().Value(userContextKey).(*store.User) //** access value from any ctx with .value method and.(return value type must be this**)
	if !ok {
		panic("missing user in request")
	}
	return user
}

// this func would authenticate incoming req by attaching token on the req with context & if it does --> calls the next fnc 
func (usrmw *UserMiddleware) Authenticate(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter,r *http.Request){

		// add header to r with type "vary"  and "auth" data
		w.Header().Set("vary","authorization")

		// extract header from *req
		extractedHeader := r.Header.Get("authorization")
		
		// validate if there is header coming or not ~ deal w/ it
		if extractedHeader == "" {
			usrmw.Logger.Warn("unexpected error occurred","error","user must be logged in")
			// set Header with user context attached via fnc as anon user with missing user ctx
			r = SetUser(store.AnonyUser,r) //* label as anony user if user is missing from ctx
			//must call next to fulfil return type and prupose of mw fnc
			next.ServeHTTP(w,r)
			return
		}

		splitedHeader := strings.Split(extractedHeader," ")
		if len(splitedHeader) != 2 || splitedHeader[0] != "Bearer" {
			usrmw.Logger.Error("unexpected error occurred","error","invalid headers")
			// err handeling
			utils.WriteJson(w,http.StatusUnauthorized,utils.Envelop{"error":"invalid authorization header"})
			return
		}
		// if it exists --> split it
		// extract from it via indexing
		token := splitedHeader[1]
		
		//  get user token 
		user,err := usrmw.UserStore.GetUserToken(tokens.ScopeAuthForTokensScope,token) //* finds that entry which is in users and token table with user_id ref id
		// multi different types of err checks
		if err != nil {
			usrmw.Logger.Error("unexpected error occurred","error",err)
			utils.WriteJson(w,http.StatusUnauthorized,utils.Envelop{"error":"invalid token or expired"})
			return
		}
		if user == nil {
			usrmw.Logger.Error("unexpected error occurred","error","user's token is either expired or invalid")
			utils.WriteJson(w,http.StatusUnauthorized,utils.Envelop{"error":"invalid token or expired"})
			return
		}
		//  if everything goes wring set user with context attached ✅✅
		r = SetUser(user,r)
		// call next to serve purpose of mw fnc
		next.ServeHTTP(w,r)
	})
}



// func that serves as bg to check if requesting user has token in ctx or not --> otherwise won't call the next route fnc
func (usrmw *UserMiddleware) RequiresAuthorization(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter,r *http.Request){
		// check for "user" in req context
		user := GetUser(r)

		// checks if user is anonymouse or not
		if user.CheckUser() {
			usrmw.Logger.Error("unexpected error occurred","error","user must be logged in")
			// if user is anon as user is not returned from getUser fnc that retireves user from its ctx
			utils.WriteJson(w,http.StatusUnauthorized,utils.Envelop{"error" : "you must be logged in to view this page!."})
			return
		}

		//  if user is authorized and has user in ctx --> serve the next route ✅✅
		next.ServeHTTP(w,r)
	})
}