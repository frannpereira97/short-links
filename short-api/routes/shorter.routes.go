package routes

import (
	"encoding/json"
	"net/http"
	"os"
	"time"

	"short-api/database"
	"short-api/models"

	"github.com/asaskevich/govalidator"
	"github.com/go-co-op/gocron"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"gorm.io/gorm"
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
	var user models.User
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
			tokenH := r.Header.Get("x-jwt-token")
			claims := GetClaims(tokenH)

			username := claims["usuario"]

			database.DB.Where("user_name = ?", username).First(&user)

			short.UserID = user.ID
			//Crea el short con el short deseado
			createdShort := database.DB.Create(&short)
			err := createdShort.Error
			if err != nil {
				w.WriteHeader(http.StatusBadRequest) // 400
				w.Write([]byte(err.Error()))
			}
			// Agrego el dominio al short para la que la respuesta sea completa
			short.Short = os.Getenv("DOMAIN") + "/" + short.Short
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
		tokenH := r.Header.Get("x-jwt-token")
		claims := GetClaims(tokenH)

		username := claims["usuario"]

		database.DB.Where("user_name = ?", username).First(&user)

		//TESTING

		newShort := uuid.New().String()[:6]
		short.Short = newShort

		short.UserID = user.ID

		createdShort := database.DB.Create(&short)

		err := createdShort.Error
		if err != nil {
			w.WriteHeader(http.StatusBadRequest) // 400
			w.Write([]byte(err.Error()))
		}
		short.Short = os.Getenv("DOMAIN") + "/" + short.Short
		json.NewEncoder(w).Encode(&short)

		return

	}

}

type ShortJSON struct {
	ID     uint   `json:"id"`
	Short  string `json:"short"`
	Pagina string `json:"pagina"`
}

func GetShortsHandler(w http.ResponseWriter, r *http.Request) {
	var shorts []models.Short

	tokenH := r.Header.Get("x-jwt-token")
	claims := GetClaims(tokenH)
	userId := GetUserID(claims["usuario"].(string))

	database.DB.Where("permisos = ? OR user_id = ?", "public", userId).Find(&shorts)

	var shortsJSON []ShortJSON
	for _, s := range shorts {

		shortsJSON = append(shortsJSON, ShortJSON{
			ID:     s.ID,
			Short:  os.Getenv("DOMAIN") + "/" + s.Short,
			Pagina: s.Pagina,
		})
	}

	jsonResult, err := json.Marshal(shortsJSON)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonResult)
}

type ShortEditJSON struct {
	Short    string `json:"short"`
	NewShort string `json:"newShort"`
}

func EditShortHandler(w http.ResponseWriter, r *http.Request) {
	var short models.Short
	var shortEdit ShortEditJSON

	json.NewDecoder(r.Body).Decode(&shortEdit)
	database.DB.Where("short = ?", shortEdit.Short).First(&short)
	if short.ID == 0 {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("Short no encontrado"))
		w.Write([]byte(short.Short))
		return
	}
	short.Short = shortEdit.NewShort
	database.DB.Save(&short)

	json.NewEncoder(w).Encode(&short)
}

func DeleteShortHandler(w http.ResponseWriter, r *http.Request) {
	var short models.Short
	var shortEdit ShortEditJSON
	json.NewDecoder(r.Body).Decode(&shortEdit)

	database.DB.Where("short = ?", shortEdit.Short).First(&short)
	if short.ID == 0 {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("Short no encontrado"))
		w.Write([]byte(short.Short))
		return
	}
	database.DB.Delete(&short)
	w.Write([]byte("Short eliminado"))
}

func ExpiredShorts(db *gorm.DB) error {
	var shorts models.Short
	result := db.Where("expiry < ?", time.Now()).Delete(&shorts)
	if result.Error != nil {

		return result.Error
	}
	// Opcional: Puedes usar result.RowsAffected para saber cuántas filas se han eliminado
	return nil
}

func CronJobs() {
	s := gocron.NewScheduler(time.UTC)

	// 4
	s.Every(1).Days().Do(ExpiredShorts, database.DB)

	// 5
	s.StartBlocking()
}
