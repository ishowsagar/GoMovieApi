# Curl requests (Ubuntu)

BASE_URL="http://localhost:8080"

## Health

curl -i "http://localhost:8080/health"

## Register user

curl -i -X POST "http://localhost:8080/api/users/register" \
 -H "Content-Type: application/json" \
 -d '{
"username":"testuser2",
"password":"Pass@123",
"email":"testuser2@example.com",
"bio":"test profile"
}'

## Login (create token)

curl -i -X POST "http://localhost:8080/api/tokens/authentication" \
 -H "Content-Type: application/json" \
 -d '{
"username":"testuser2",
"password":"Pass@123"
}'

## Save token to variable (needs jq)

TOKEN=$(curl -s -X POST "http://localhost:8080/api/tokens/authentication" \
 -H "Content-Type: application/json" \
 -d '{"username":"testuser1","password":"Pass@123"}' | jq -r '.auth_token.token')
echo "$TOKEN"

## Protected route without token (should be 401)

curl -i "http://localhost:8080/api/movies/all"

## Get all movies (with token)

curl -i "http://localhost:8080/api/movies/all" \
 -H "Authorization: Bearer GNZQFB2AYTAY6WC4BXJZR5J7OSI5CJ3EZGMO7LDI5KLIGHNGCQMQ"

## Create movie

curl -i -X POST "http://localhost:8080/api/movie/create" \
 -H "Authorization: Bearer GNZQFB2AYTAY6WC4BXJZR5J7OSI5CJ3EZGMO7LDI5KLIGHNGCQMQ" \
 -H "Content-Type: application/json" \
 -d '{
"name":"Inception",
"genre":"Sci-fi",
"description":"Dream layers and time distortion",
"ratings":8.8
}'

## Get movie by id

MOVIE_ID="replace-with-real-id"
curl -i "http://localhost:8080/api/movies/movie/$MOVIE_ID" \
 -H "Authorization: Bearer $TOKEN"

## Update movie by id

curl -i -X PUT "http://localhost:8080/api/movies/movie/update/$MOVIE_ID" \
 -H "Authorization: Bearer $TOKEN" \
 -H "Content-Type: application/json" \
 -d '{
"name":"Inception Updated",
"genre":"Sci-fi",
"description":"Updated description",
"ratings":9.0
}'

## Delete movie by id

curl -i -X DELETE "http://localhost:8080/api/movies/movie/delete/$MOVIE_ID" \
 -H "Authorization: Bearer $TOKEN"

## Delete all movies

curl -i -X DELETE "http://localhost:8080/api/movies/all/delete" \
 -H "Authorization: Bearer $TOKEN"

## Negative auth tests

curl -i "http://localhost:8080/api/movies/all" -H "Authorization: GNZQFB2AYTAY6WC4BXJZR5J7OSI5CJ3EZGMO7LDI5KLIGHNGCQMQ"
curl -i "http://localhost:8080/api/movies/all" -H "Authorization: Bearer invalidtoken"

# Curl requests for testing API

# Post - api/movies/create

curl -i -X POST "http://localhost:8080/api/movie/create" \
 -H "Content-Type: application/json" \
 -d '{
"name" : "Marvel Avengers Endgame",
"genre" : "Sci-fic",
"description" : "Thanos has already destroyed Planets by snapping his finger but wait, There is something big awaiting to take place",
"ratings" : 8.9
}'

# APP FLOW

Cmd/server/main.go --> sets up server --> passing req to handler via Chi router in router.go
Router.go --> routed requests on routes invokes functions on controllers functions --> sends response to client
Controllers func use methods belongs to Movie type --> all methods query db calls
Db connection is made in db.go with sql.open(passin in driver pgx from jackc, and db connection string) --> which stores that db connectioin in a instance of type Db which stores \*sql.Db type of connection <---> which is returned by db.go connectToDb function
Api --> stores Movie type in a type through which a function returns empty instance of that type but that fnc takes in a db parameter and assigns that db connection arg to var which stores same type of function in api package
that var now stores db connection { note : it stores db connection actual not concrete type that stores db connection} --> used by services~dbCalls methods to query db calls
and methods passed to controller functions and those functions are that fnc which are used by the router

update :
Added Interface to re-route api db calls through MovieMethodStore Interface --> in services.Go --> created interface that stores all the methods that belongs to type Movie which essentially do all the db calls and without exposing it directly --> making it "dependency injection" implemented
--> interface stores all methods which make db calls --> interface is implemented by the type Movie --> so In controllers.go we made a type to store that interface and made a function to intiate the instance of that interface to meet the type that implements it --> passed type that stores intrface in the method reciever in all the controllers functions to turn them into methods --> now all the controllers methods belongs to type (m MovieMethodStore) and could be accessed via making a instance of type --> as they belongs to type 'm' --> accessed through var var of type or instance

db calls rules with context :
Most-used methods on DB:

QueryContext(ctx, query, args...)
Use when query returns multiple rows.
Returns rows object, then loop with rows.Next() and rows.Scan(...).

QueryRowContext(ctx, query, args...)
Use when query returns exactly one row.
Returns single row object, then row.Scan(...).

ExecContext(ctx, query, args...)
Use when query does not return row data.
Examples: INSERT (without returning), UPDATE, DELETE, CREATE TABLE.
Returns result, then check RowsAffected or LastInsertId (driver-dependent).

PrepareContext(ctx, query)
Use when same statement runs many times.
Returns prepared statement, then stmt.QueryContext / stmt.ExecContext.


# next step

    Adding Auth --> JWT token
        Serve no route if token not found

# workflow

1. Create a method in services.md --> must add (m Movie) reciever so that it belongs to type and accessible through it
2. create corresponding method in controllers.go to handle that services.go method call --> must add (m MovieMethodStore) in reciever so it belongs to it and accessible in router.go
3. Add route for the controller method and wire it up
   ⭐Remember ==> Don't forget to add method signature in the interface that stores all the method from type "Movie" that implements it


```
