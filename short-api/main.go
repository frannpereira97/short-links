package main

import (
	"log"
	"net/http"
	"os"
	"time"

	"short-api/database"
	"short-api/models"
	"short-api/routes"

	"github.com/didip/tollbooth/v7"
	"github.com/didip/tollbooth/v7/limiter"
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

	//Ejecuto el cron

	//Cargo el router
	r := mux.NewRouter()

	//Cargo los archivos para que corra el front
	fs := http.FileServer(http.Dir("./web/assets"))
	r.PathPrefix("/assets/").Handler(http.StripPrefix("/assets/", fs))

	//Limitador de solicitudes
	limiter := tollbooth.NewLimiter(2, &limiter.ExpirableOptions{DefaultExpirationTTL: time.Second})
	limiter.SetMessage("Has alcanzado el limite de solicitudes.")
	limiter.SetMessageContentType("application/json; charset=utf-8")

	//Login/Logout Request
	r.Handle("/users/login", tollbooth.LimitFuncHandler(limiter, routes.LoginHandler)).Methods("POST")
	r.Handle("/users/validate", tollbooth.LimitFuncHandler(limiter, routes.ValidateLoginHandler)).Methods("POST")
	r.Handle("/users/logout", tollbooth.LimitFuncHandler(limiter, routes.LogoutHandler)).Methods("POST")

	//Redirigir a la direccion del Short
	r.Handle("/{url}", tollbooth.LimitFuncHandler(limiter, routes.ResolveURL)).Methods("GET")

	//Usuarios
	r.Handle("/users/changedata", tollbooth.LimitFuncHandler(limiter, routes.WithJWTAuth(routes.ChangeDataHandler))).Methods("POST")
	r.Handle("/users/changepw", tollbooth.LimitFuncHandler(limiter, routes.WithJWTAuth(routes.ChangePasswordHandler))).Methods("POST")
	r.Handle("/users/list", tollbooth.LimitFuncHandler(limiter, routes.WithJWTAuth(routes.GetUsersHandler))).Methods("GET")
	r.Handle("/users/data", tollbooth.LimitFuncHandler(limiter, routes.WithJWTAuth(routes.GetUserDataHandler))).Methods("GET")

	r.Handle("/users/create", tollbooth.LimitFuncHandler(limiter, routes.CreateUserHandler)).Methods("POST")

	//Revisar
	r.Handle("/users/{id}", tollbooth.LimitFuncHandler(limiter, routes.WithJWTAuth(routes.DeleteUsersHandler))).Methods("DELETE")
	r.Handle("/users/{id}", tollbooth.LimitFuncHandler(limiter, routes.WithJWTAuth(routes.GetUserHandler))).Methods("GET")
	//Crear Short
	r.Handle("/users/Shorten", tollbooth.LimitFuncHandler(limiter, routes.WithJWTAuth(routes.ShortenURL))).Methods("POST")

	//Listar todos los Short y enviarlos
	r.Handle("/shorts/list", tollbooth.LimitFuncHandler(limiter, routes.WithJWTAuth(routes.GetShortsHandler))).Methods("GET")
	r.Handle("/shorts/edit", tollbooth.LimitFuncHandler(limiter, routes.WithJWTAuth(routes.EditShortHandler))).Methods("POST")
	r.Handle("/shorts/delete", tollbooth.LimitFuncHandler(limiter, routes.WithJWTAuth(routes.DeleteShortHandler))).Methods("POST")

	//CORS
	headersOk := handlers.AllowedHeaders([]string{"X-Requested-With", "Content-Type", "x-jwt-token"})
	originsOk := handlers.AllowedOrigins([]string{os.Getenv("ORIGIN_ALLOWED")})
	methodsOk := handlers.AllowedMethods([]string{"GET", "HEAD", "POST", "PUT", "OPTIONS"})

	http.ListenAndServe(":4000", handlers.CORS(originsOk, headersOk, methodsOk)(r))

	routes.CronJobs()

}
