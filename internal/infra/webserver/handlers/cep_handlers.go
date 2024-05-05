package handlers

import (
	"encoding/json"
	"fmt"
	"io"
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

	externalUrl := fmt.Sprintf("http://localhost:8081/temperatura/%s", cepInput.Cep)
	fmt.Printf("external URL: %s\n", externalUrl)
	resp, err := http.Get(externalUrl)
	if err != nil {
		panic(err)
	}
	w.WriteHeader(resp.StatusCode)
	defer resp.Body.Close()
	// w.Header().Set("Content-Type", "application/json")
	io.Copy(w, resp.Body)
	return
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
