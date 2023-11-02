package routes

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/frannpereira97/short-links/database"
	"github.com/frannpereira97/short-links/models"
	"github.com/gorilla/mux"
	"golang.org/x/crypto/bcrypt"
)

func GetUsersHandler(w http.ResponseWriter, r *http.Request) {
	var users []models.User
	database.DB.Find(&users)
	json.NewEncoder(w).Encode(&users)
}

func GetUserHandler(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	var user models.User
	database.DB.First(&user, params["id"])

	if user.ID == 0 {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("Usuario no encontrado"))
		return

	}
	json.NewEncoder(w).Encode(&user)

}

type userData struct {
	Nombre       string `json:"nombre"`
	Apellido     string `json:"apellido"`
	Usuario      string `json:"UserName"`
	Contraseña   string `json:"Password"`
	Email        string `json:"email"`
	Sexo         string `json:"sexo"`
	Nacionalidad string `json:"nacionalidad"`
	Provincia    string `json:"provincia"`
	Ciudad       string `json:"ciudad"`
	Domicilio    string `json:"Domicilio"`
}

func CreateUserHandler(w http.ResponseWriter, r *http.Request) {
	var user models.User
	var datos models.Datos
	var uData userData

	json.NewDecoder(r.Body).Decode(&uData)

	//Creo Usuario
	user.UserName = uData.Usuario
	user.Password = uData.Contraseña
	user.Email = uData.Email

	tokenString, err2 := createJWT(&user)
	if err2 != nil {
		w.WriteHeader(http.StatusBadRequest) // 400
		w.Write([]byte(err2.Error()))
	}
	user.Token = tokenString
	user.Password, err2 = hashPW(user.Password)
	if err2 != nil {
		w.WriteHeader(http.StatusBadRequest) // 400
		w.Write([]byte(err2.Error()))
	}
	createdUser := database.DB.Create(&user)
	err := createdUser.Error
	if err != nil {
		w.WriteHeader(http.StatusBadRequest) // 400
		w.Write([]byte(err.Error()))
		fmt.Println(err.Error())
	}

	//Cargo Datos
	id := GetUserID(user.UserName)

	datos.Nombre = uData.Nombre
	datos.Apellido = uData.Apellido
	datos.Sexo = uData.Sexo
	datos.Nacionalidad = uData.Nacionalidad
	datos.Provincia = uData.Provincia
	datos.Ciudad = uData.Ciudad
	datos.Domicilio = uData.Domicilio
	datos.UserID = uint(id)

	userData := database.DB.Create(&datos)
	err2 = userData.Error
	if err2 != nil {
		w.WriteHeader(http.StatusBadRequest) // 400
		w.Write([]byte(err2.Error()))
		fmt.Println(err2.Error())
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{
		"message":     "Login correcto",
		"x-jwt-token": tokenString,
		"redirectTo":  "/home",
	})
}

func hashPW(password string) (string, error) {

	encpw, err3 := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)

	if err3 != nil {
		return "no funciona", err3
	}
	return string(encpw), nil
}

func GetUserID(username string) int {

	var userID models.User

	database.DB.Where("user_name = ?", username).First(&userID)

	if userID.ID == 0 {
		return -1
	}

	id := userID.ID

	return int(id)

}

func DeleteUsersHandler(w http.ResponseWriter, r *http.Request) {
	var user models.User
	params := mux.Vars(r)

	database.DB.First(&user, params["id"])

	if user.ID == 0 {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("Usuario no encontrado para eliminar"))
		return
	}
	database.DB.Unscoped().Delete(&user)

}
