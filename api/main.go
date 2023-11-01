package main

import (
	"log"
	"net/http"
	"os"
	"time"

	"github.com/didip/tollbooth/v7"
	"github.com/didip/tollbooth/v7/limiter"
	"github.com/frannpereira97/short-links/database"
	"github.com/frannpereira97/short-links/models"
	"github.com/frannpereira97/short-links/routes"
	jwt "github.com/golang-jwt/jwt/v4"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
)

type MyCustomClaims struct {
	Foo string `json:"foo"`
	jwt.RegisteredClaims
}

func main() {
	//Cargo las variables de entorno
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	database.Dbconnect()
	//Cargo los modelos en la BD
	database.DB.AutoMigrate(models.Short{})
	database.DB.AutoMigrate(models.Datos{})
	database.DB.AutoMigrate(models.User{})

	r := mux.NewRouter()

	//Cargo los archivos para que corra el front
	fs := http.FileServer(http.Dir("./web/assets"))
	r.PathPrefix("/assets/").Handler(http.StripPrefix("/assets/", fs))

	//Limitador de solicitudes
	limiter := tollbooth.NewLimiter(1, &limiter.ExpirableOptions{DefaultExpirationTTL: time.Second})
	limiter.SetMessage("Has alcanzado el limite de solicitudes.")
	limiter.SetMessageContentType("application/json; charset=utf-8")

	//Rutas
	// r.HandleFunc("/hello", routes.HomeHandler)

	// Index y Login
	r.Handle("/", tollbooth.LimitFuncHandler(limiter, routes.IndexHandler)).Methods("GET")
	r.Handle("/home", tollbooth.LimitFuncHandler(limiter, routes.HomeHandler)).Methods("GET")
	r.Handle("/users/login", tollbooth.LimitFuncHandler(limiter, routes.LoginHandler)).Methods("POST")
	r.Handle("/users/logout", tollbooth.LimitFuncHandler(limiter, routes.LogoutHandler)).Methods("POST")

	//Usuarios
	r.Handle("/users", tollbooth.LimitFuncHandler(limiter, routes.GetUsersHandler)).Methods("GET")
	r.Handle("/users/{id}", tollbooth.LimitFuncHandler(limiter, routes.WithJWTAuth(routes.GetUserHandler))).Methods("GET")
	r.Handle("/users", tollbooth.LimitFuncHandler(limiter, routes.CreateUserHandler)).Methods("POST")
	r.Handle("/users/{id}", tollbooth.LimitFuncHandler(limiter, routes.DeleteUsersHandler)).Methods("DELETE")

	//Shorts
	r.Handle("/users/Shorten", tollbooth.LimitFuncHandler(limiter, routes.WithJWTAuth(routes.ShortenURL))).Methods("POST")
	r.Handle("/{url}", tollbooth.LimitFuncHandler(limiter, routes.ResolveURL)).Methods("GET")

	//CORS
	headersOk := handlers.AllowedHeaders([]string{"X-Requested-With", "Content-Type", "x-jwt-token"})
	originsOk := handlers.AllowedOrigins([]string{os.Getenv("ORIGIN_ALLOWED")})
	methodsOk := handlers.AllowedMethods([]string{"GET", "HEAD", "POST", "PUT", "OPTIONS"})

	http.ListenAndServe(":4000", handlers.CORS(originsOk, headersOk, methodsOk)(r))

}
