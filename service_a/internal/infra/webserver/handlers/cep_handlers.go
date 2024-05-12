package handlers

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"regexp"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/propagation"
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

	tracer := otel.Tracer("cep-validator-tracer")
	ctx, span := tracer.Start(ctx, "cep-temperatura-request")
	defer span.End()

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
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, externalUrl, nil)
	if err != nil {
		http.Error(
			w, "erro ao preparar requisição para buscar temperatura",
			http.StatusInternalServerError,
		)
		return
	}
	otel.GetTextMapPropagator().Inject(ctx, propagation.HeaderCarrier(req.Header))
	fmt.Printf("external URL: %s\n", externalUrl)
	resp, err := http.DefaultClient.Do(req)
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
