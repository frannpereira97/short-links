package routes

import (
	"fmt"
	"net/http"
	"os"

	"github.com/frannpereira97/short-links/database"
	"github.com/frannpereira97/short-links/models"
	jwt "github.com/golang-jwt/jwt/v4"
	"github.com/gorilla/mux"
)

func createJWT(user *models.User) (string, error) {

	signingKey := os.Getenv("JWT_SECRET")

	claims := &jwt.MapClaims{
		"ExpiresAt": 15000,
		"usuario":   user.UserName,
		"permisos":  user.Permisos,
	}

	fmt.Println(claims)
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString([]byte(signingKey))
}

func WithJWTAuth(handlerFunc http.HandlerFunc) http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {
		tokenString := r.Header.Get("x-jwt-token")
		if tokenString == "" {
			//No tiene token
			w.WriteHeader(http.StatusForbidden)
			w.Write([]byte("No tenes permisos para acceder"))
			return
		}

		token, err := validateJWT(tokenString)
		if err != nil {
			w.WriteHeader(http.StatusNotFound)
			w.Write([]byte("Token invalido"))
			return
		}
		claims := token.Claims.(jwt.MapClaims)

		//TESTING
		params := mux.Vars(r)
		var user models.User
		database.DB.First(&user, params["id"])

		fmt.Println(claims["usuario"], user.UserName)
		if claims["usuario"] != user.UserName {
			w.WriteHeader(http.StatusNotFound)
			w.Write([]byte("No tenes permisos para acceder"))
			return
		}

		handlerFunc(w, r)
	}
}

func validateJWT(tokenString string) (*jwt.Token, error) {

	signingKey := []byte(os.Getenv("JWT_SECRET"))

	return jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return signingKey, nil
	})

}
