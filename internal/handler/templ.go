package handler

import (
	"net/http"

	"github.com/sebasukodo/chirpy/templates"
)

func (cfg *ApiConfig) Login(w http.ResponseWriter, r *http.Request) {

	if err := templates.Login("Chirpy Login").Render(r.Context(), w); err != nil {
		respondWithError(w, 500, "Error")
		return
	}

}

func (cfg *ApiConfig) Register(w http.ResponseWriter, r *http.Request) {

	if err := templates.Register("Chirpy Register").Render(r.Context(), w); err != nil {
		respondWithError(w, 500, "Error")
		return
	}
}

func (cfg *ApiConfig) Homepage(w http.ResponseWriter, r *http.Request) {

	c := templates.HomepageContent()
	if err := templates.Layout(c, "Chirpy - short messages").Render(r.Context(), w); err != nil {
		respondWithError(w, 500, "Error")
		return
	}
}
