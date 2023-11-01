package routes

import (
	"encoding/json"
	"net/http"

	"github.com/asaskevich/govalidator"
	"github.com/frannpereira97/short-links/database"
	"github.com/frannpereira97/short-links/models"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

func ResolveURL(w http.ResponseWriter, r *http.Request) {

	var short models.Short
	vars := mux.Vars(r)

	json.NewDecoder(r.Body).Decode(&short)

	short.Short = vars["url"]

	database.DB.Where("short = ?", short.Short).First(&short)
	if short.ID == 0 {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("Short no encontrado"))
		w.Write([]byte(short.Short))
		return
	}

	value := short.Pagina
	go database.DB.Exec(`UPDATE shorts SET abierto = abierto + 1 WHERE pagina = ?`, value)

	http.Redirect(w, r, value, http.StatusMovedPermanently)

}

func ShortenURL(w http.ResponseWriter, r *http.Request) {

	var short models.Short
	//customShort := vars["customShort"]

	json.NewDecoder(r.Body).Decode(&short)

	//Verifica que la url sea valida
	if !govalidator.IsURL(short.Pagina) {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("URL invalida"))
		w.Write([]byte(short.Pagina))
		return
	}
	if short.Short != "" {
		//Verifica que el customShort no exista
		database.DB.Where("short = ?", short.Short).First(&short)
		if short.ID == 0 {
			//Crea el short con el short deseado
			createdShort := database.DB.Create(&short)
			err := createdShort.Error
			if err != nil {
				w.WriteHeader(http.StatusBadRequest) // 400
				w.Write([]byte(err.Error()))
			}
			json.NewEncoder(w).Encode(&short)

			return
		} else {
			//Avisa que ya existe
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(map[string]string{
				"message": "El short ya existe",
			})
		}

	} else if short.Short == "" {

		newShort := uuid.New().String()[:6]
		short.Short = newShort

		createdShort := database.DB.Create(&short)
		err := createdShort.Error
		if err != nil {
			w.WriteHeader(http.StatusBadRequest) // 400
			w.Write([]byte(err.Error()))
		}
		json.NewEncoder(w).Encode(&short)

		return

	}

}
