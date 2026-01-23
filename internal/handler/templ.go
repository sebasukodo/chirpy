package handler

import (
	"net/http"

	"github.com/sebasukodo/chirpy/templates"
)

func (cfg *ApiConfig) Login(w http.ResponseWriter, r *http.Request) {

	if err := templates.Login("Chirpy Login").Render(r.Context(), w); err != nil {
		respondWithError(w, r, 500, "Error")
		return
	}

}

func (cfg *ApiConfig) Register(w http.ResponseWriter, r *http.Request) {

	if err := templates.Register("Chirpy Register").Render(r.Context(), w); err != nil {
		respondWithError(w, r, 500, "Error")
		return
	}
}

func (cfg *ApiConfig) ProfilePage(w http.ResponseWriter, r *http.Request) {

	if err := templates.ProfilePage().Render(r.Context(), w); err != nil {
		respondWithError(w, r, 500, "Error")
		return
	}
}
