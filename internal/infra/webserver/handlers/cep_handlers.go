package handlers

import (
	"encoding/json"
	"net/http"
	"regexp"
)

type CepHandler struct{}

type CepInput struct {
	Cep string `json:"cep"`
}

func NewCepHandler() *CepHandler {
	return &CepHandler{}
}

func (h *CepHandler) PostCep(w http.ResponseWriter, r *http.Request) {
	var cepInput CepInput
	err := json.NewDecoder(r.Body).Decode(&cepInput)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	invalidZipcodeMessage := "Invalid zipcode"
	isCepValid := isCepValid(cepInput.Cep)
	if !isCepValid {
		w.WriteHeader(http.StatusUnprocessableEntity)
		err := json.NewEncoder(w).Encode(&invalidZipcodeMessage)
		if err != nil {
			panic(err)
		}
		return
	}
	w.WriteHeader(http.StatusCreated)
}

func isCepValid(cep string) bool {
	var re = regexp.MustCompile(`^[0-9]+$`)
	if len(cep) != 8 {
		return false
	}
	if !re.MatchString(cep) {
		return false
	}
	return true
}
