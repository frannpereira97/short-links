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

type ShortUser struct {
	ID       uint   `json:"id"`
	Username string `json:"username"`
	Permisos string `json:"permisos"`
	Admin    bool   `json:"admin"`
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

type PermisosResponse struct {
	Permisos string `json:"permisos"`
	Admin    bool   `json:"admin"`
}

func GetUserHandler(w http.ResponseWriter, r *http.Request) {
	var users []models.User
	database.DB.Find(&users)
	json.NewEncoder(w).Encode(&users)
}

func GetUsersHandler(w http.ResponseWriter, r *http.Request) {
	var user []models.User
	var shorterUser []ShortUser

	tokenH := r.Header.Get("x-jwt-token")
	claims := GetClaims(tokenH)
	permisos := claims["permisos"].(string)
	if permisos == "admin" {
		database.DB.Find(&user)
		for _, s := range user {

			shorterUser = append(shorterUser, ShortUser{
				ID:       s.ID,
				Username: s.UserName,
				Permisos: s.Permisos,
				Admin:    true,
			})
		}
		jsonResult, err := json.Marshal(shorterUser)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write(jsonResult)
	} else if permisos == "user" {
		permisos := PermisosResponse{
			Permisos: "user",
			Admin:    false,
		}
		jsonPermisos, err := json.Marshal(permisos)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write(jsonPermisos)
	} else {
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte("No tiene permisos para acceder a esta pagina"))
	}
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
	user.Permisos = "user"

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
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{
		"message":     "Login correcto",
		"x-jwt-token": tokenString,
		"redirectTo":  "/about",
	})
}

type PassEdit struct {
	Password string `json:"password"`
	NewPass  string `json:"newPass"`
	VnewPass string `json:"vnewPass"`
}

func ChangePasswordHandler(w http.ResponseWriter, r *http.Request) {
	var passEdit PassEdit
	var user models.User
	token := r.Header.Get("x-jwt-token")
	claims := GetClaims(token)
	fmt.Println(claims)
	database.DB.Where("user_name = ?", claims["usuario"]).First(&user)
	if user.ID == 0 {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("Usuario no encontrado"))
		return
	}
	json.NewDecoder(r.Body).Decode(&passEdit)
	chk := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(passEdit.Password))
	if chk != nil {
		fmt.Println("Contraseña incorrecta")
		w.WriteHeader(http.StatusBadRequest)
		http.Error(w, "Error", http.StatusBadRequest)
		return
	} else {
		if passEdit.NewPass != passEdit.VnewPass {
			w.WriteHeader(http.StatusBadRequest)
			http.Error(w, "Error", http.StatusBadRequest)
			return
		} else {

			newPW, err2 := hashPW(passEdit.NewPass)
			if err2 != nil {
				w.WriteHeader(http.StatusBadRequest) // 400
				w.Write([]byte(err2.Error()))
			}
			user.Password = newPW
			database.DB.Model(&user).Where("user_name = ?", user.UserName).Update("password", user.Password)
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(map[string]string{
				"message": "Contraseña cambiada correctamente",
			})
		}
	}
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

func GetUserDataHandler(w http.ResponseWriter, r *http.Request) {
	var datos models.Datos
	var user models.User
	var uData userData

	token := r.Header.Get("x-jwt-token")
	claims := GetClaims(token)
	id := GetUserID(claims["usuario"].(string))
	database.DB.Where("user_id = ?", id).First(&datos)
	if datos.ID == 0 {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("Usuario no encontrado"))
		return
	}
	database.DB.Where("id = ?", id).First(&user)
	if user.ID == 0 {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("Usuario no encontrado"))
		return
	}
	uData.Nombre = datos.Nombre
	uData.Apellido = datos.Apellido
	uData.Usuario = user.UserName
	uData.Email = user.Email
	uData.Sexo = datos.Sexo
	uData.Nacionalidad = datos.Nacionalidad
	uData.Provincia = datos.Provincia
	uData.Ciudad = datos.Ciudad
	uData.Domicilio = datos.Domicilio
	fmt.Println(uData)
	json.NewEncoder(w).Encode(&uData)
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
