package routes

import (
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"

	"github.com/frannpereira97/short-links/database"
	"github.com/frannpereira97/short-links/models"
	"golang.org/x/crypto/bcrypt"
)

type LoginRequest struct {
	UserName string
	Password string
}

var tmpl = template.Must(template.ParseGlob("web/*.html"))

func RegisterHandler(w http.ResponseWriter, r *http.Request) {
	tmpl.ExecuteTemplate(w, "register.html", nil)
}

func HomeHandler(w http.ResponseWriter, r *http.Request) {
	tmpl.ExecuteTemplate(w, "home.html", nil)
}

func IndexHandler(w http.ResponseWriter, r *http.Request) {
	tmpl.ExecuteTemplate(w, "login.html", nil)
}
func AboutHandler(w http.ResponseWriter, r *http.Request) {

	tmpl.ExecuteTemplate(w, "about.html", nil)
}

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	var login LoginRequest
	//Almaceno en JSON los datos de log
	if err := json.NewDecoder(r.Body).Decode(&login); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}
	//Obtengo el ID
	username := GetUserID(login.UserName)
	//Si no existe el usuario
	if username == -1 {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("Usuario no encontrado"))
		return
	} else if username != -1 {
		//Obtengo la contraseña
		var user models.User
		database.DB.Where("id = ?", username).First(&user)
		//Reviso la contraseña y la comparo
		chk := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(login.Password))
		if chk != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("Contraseña incorrecta"))
			return
		} else {
			//Creo la sesion - Asigno token
			tokenString, err2 := createJWT(&user)
			if err2 != nil {
				w.WriteHeader(http.StatusBadRequest) // 400
				w.Write([]byte(err2.Error()))
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
	fmt.Println("Logout correcto")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{
		"message":    "Logout correcto",
		"redirectTo": "/",
	})
}
