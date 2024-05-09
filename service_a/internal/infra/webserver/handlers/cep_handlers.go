package handlers

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"regexp"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/metric"
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

func init() {
	var err error
	viewCounter, err = meter.Int64Counter("user.views",
		metric.WithDescription("The number of views"),
		metric.WithUnit("{views}"))
	if err != nil {
		panic(err)
	}
}

func (h *CepHandler) PostCep(w http.ResponseWriter, r *http.Request) {
	ctx, span := tracer.Start(r.Context(), "cep-validation")
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
