package handlers

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"regexp"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/propagation"
)

var (
	tracer      = otel.Tracer("cep-validation-request")
	meter       = otel.Meter("cep-validation-info-service")
	viewCounter metric.Int64Counter
)

type CepHandler struct{}

type CepInput struct {
	Cep string `json:"cep"`
}

func NewCepHandler() *CepHandler {
	return &CepHandler{}
}

func (h *CepHandler) PostCep(w http.ResponseWriter, r *http.Request) {
	carrier := propagation.HeaderCarrier(r.Header)
	ctx := r.Context()
	ctx = otel.GetTextMapPropagator().Extract(ctx, carrier)

	tracer := otel.Tracer("cep-validator-request-tracer")
	ctx, span := tracer.Start(ctx, "cep-validator-temperatura-request")

	defer span.End()

	viewCounter.Add(ctx, 1)

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

	externalUrl := fmt.Sprintf("http://temperatura-cep:8081/temperatura/%s", cepInput.Cep)
	fmt.Printf("external URL: %s\n", externalUrl)
	resp, err := http.Get(externalUrl)
	if err != nil {
		http.Error(w, "erro ao buscar temperatura", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(resp.StatusCode)
	defer resp.Body.Close()
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
