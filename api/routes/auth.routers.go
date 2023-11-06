package routes

import (
	"net/http"
	"os"

	"github.com/frannpereira97/short-links/database"
	"github.com/frannpereira97/short-links/models"
	jwt "github.com/golang-jwt/jwt/v4"
)

func createJWT(user *models.User) (string, error) {

	signingKey := os.Getenv("JWT_SECRET")

	claims := &jwt.MapClaims{
		"ExpiresAt": 15000,
		"usuario":   user.UserName,
		"permisos":  user.Permisos,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString([]byte(signingKey))
}

func GetClaims(tokenH string) jwt.MapClaims {

	token, _, err2 := new(jwt.Parser).ParseUnverified(tokenH, jwt.MapClaims{})
	if err2 != nil {
		return nil
	}
	claims := token.Claims.(jwt.MapClaims)

	return claims
}

func ValidateLoginHandler(w http.ResponseWriter, r *http.Request) {
	tokenH := r.Header.Get("x-jwt-token")
	valid, err := validateJWT(tokenH)
	if err != nil || !valid.Valid {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("Token invalido"))
		return
	} else {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Token valido"))
	}
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
		var user models.User
		database.DB.Where("user_name = ?", claims["usuario"]).First(&user)
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
