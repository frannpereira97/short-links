package routes

import (
	"encoding/json"
	"net/http"

	"short-api/database"
	"short-api/models"

	"golang.org/x/crypto/bcrypt"
)

type LoginRequest struct {
	UserName string
	Password string
}

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	var login LoginRequest
	//Almaceno en JSON los datos de log
	if err := json.NewDecoder(r.Body).Decode(&login); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		http.Error(w, "Error", http.StatusBadRequest)
		return
	}
	//Obtengo el ID
	username := GetUserID(login.UserName)
	//Si no existe el usuario
	if username == -1 {
		http.Error(w, "Error", http.StatusBadRequest)
		w.WriteHeader(http.StatusBadRequest)
		return
	} else if username != -1 {
		//Obtengo la contraseña
		var user models.User
		database.DB.Where("id = ?", username).First(&user)
		//Reviso la contraseña y la comparo
		chk := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(login.Password))
		if chk != nil {
			w.WriteHeader(http.StatusBadRequest)
			http.Error(w, "Error", http.StatusBadRequest)
			return
		} else {
			//Creo la sesion - Asigno token
			tokenString, err2 := createJWT(&user)
			if err2 != nil {
				w.WriteHeader(http.StatusBadRequest) // 400
				http.Error(w, "Error", http.StatusBadRequest)
				return
			}
			//Redirecciono
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(map[string]string{
				"message":     "Login correcto",
				"x-jwt-token": tokenString,
				"redirectTo":  "/about",
			})
		}
	}
}

func LogoutHandler(w http.ResponseWriter, r *http.Request) {
	// Eliminar token
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{
		"message":    "Logout correcto",
		"redirectTo": "/",
	})
}
