package companies

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/go-chi/chi"
	"github.com/marianodsr/nura-api/middlewares"
)

func RegisterRoutes(r chi.Router) {

	r.With(middlewares.RequireAuth).Put("/{id}", updateCompanyHandler)
	r.With(middlewares.RequireAuth).Post("/", createCompanyHandler)
	r.With(middlewares.RequireAuth).Get("/{id}", getCompanyHandler)

}

func getCompanyHandler(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if id == "" {
		http.Error(w, "No id provided", http.StatusBadRequest)
		return
	}
	intID, err := strconv.Atoi(id)
	if err != nil {
		http.Error(w, "Id provided is invalid", http.StatusBadRequest)
		return
	}
	company := GetCompanyByID(uint(intID))
	json.NewEncoder(w).Encode(company)

}

func createCompanyHandler(w http.ResponseWriter, r *http.Request) {
	company := &Company{}
	if err := json.NewDecoder(r.Body).Decode(company); err != nil {
		http.Error(w, "Invalid company paramters", http.StatusBadRequest)
		return
	}
	fmt.Printf("\nFROM REQUEST: \n%+v\n", company)
	if err := CreateCompany(company); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	json.NewEncoder(w).Encode(company)
}

func updateCompanyHandler(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if id == "" {
		http.Error(w, "No id provided", 400)
		return
	}
	fmt.Printf("\nBODY: %v\n", r.Body)
	requestCompany := &Company{}
	err := json.NewDecoder(r.Body).Decode(requestCompany)
	if err != nil {
		http.Error(w, "Invalid body format", 400)
		return
	}
	fmt.Printf("\nCompany: %+v\n", requestCompany)

	intID, err := strconv.Atoi(id)
	if err != nil {
		http.Error(w, "Invalid id provided", 400)
		return
	}
	company := GetCompanyByID(uint(intID))
	company = requestCompany
	err = UpdateCompany(company)
	if err != nil {
		http.Error(w, err.Error(), 400)
	}
}
