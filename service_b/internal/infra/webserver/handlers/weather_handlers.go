package handlers

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"net/url"
	"regexp"
	"strconv"

	"github.com/go-chi/chi"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/propagation"
)

type WeatherHandler struct{}

type Localizacao struct {
	Cep         string `json:"cep"`
	Logradouro  string `json:"logradouro"`
	Complemento string `json:"complemento"`
	Bairro      string `json:"bairro"`
	Localidade  string `json:"localidade"`
	Uf          string `json:"uf"`
}

type CurrentWeather struct {
	TempC float64 `json:"temp_C"`
	TempF float64 `json:"temp_F"`
}

type WeatherApiResponse struct {
	Current CurrentWeather `json:"current"`
}

type WeatherResponse struct {
	City  string  `json:city`
	TempC float64 `json:"temp_C"`
	TempF float64 `json:"temp_F"`
	TempK float64 `json:"temp_K"`
}

func NewWeatherHandler() *WeatherHandler {
	return &WeatherHandler{}
}

func (h *WeatherHandler) GetWeather(w http.ResponseWriter, r *http.Request) {
	carrier := propagation.HeaderCarrier(r.Header)
	ctx := r.Context()
	ctx = otel.GetTextMapPropagator().Extract(ctx, carrier)

	tracer := otel.Tracer("temperatura-cep-request-tracer")
	ctx, span := tracer.Start(ctx, "temperatura-cep-request")
	defer span.End()

	invalidZipcodeMessage := "Invalid zipcode"
	zipcodeNotFound := "can not find zipcode"
	cep := chi.URLParam(r, "cep")
	isValid := isValidCep(cep)

	if !isValid {
		w.WriteHeader(http.StatusUnprocessableEntity)
		err := json.NewEncoder(w).Encode(&invalidZipcodeMessage)
		if err != nil {
			panic(err)
		}
		return
	}

	cepResponse := buscaCEP(cep)
	if cepResponse.Cep == "" {
		w.WriteHeader(http.StatusNotFound)
		err := json.NewEncoder(w).Encode(&zipcodeNotFound)
		if err != nil {
			panic(err)
		}
		return
	}
	log.Printf("CEP encontrado: %v", cepResponse)
	temperatura := buscaTemperatura(cepResponse)
	log.Printf("WeatherApiResponse: %v", temperatura)
	kelvin := temperatura.Current.TempC + 273
	response := WeatherResponse{
		City:  cepResponse.Localidade,
		TempC: temperatura.Current.TempC,
		TempF: temperatura.Current.TempF,
		TempK: kelvin,
	}
	json.NewEncoder(w).Encode(&response)
	return
}

func isValidCep(cep string) bool {
	var re = regexp.MustCompile(`^[0-9]+$`)
	if !re.MatchString(cep) {
		return false
	}
	return true
}

func buscaCEP(cep string) Localizacao {
	address := "http://viacep.com.br/ws/" + cep + "/json/"
	req, err := http.Get(address)
	if err != nil {
		panic(err)
	}
	if req.StatusCode != http.StatusOK {
		panic("Erro ao fazer requisição para ViaCEP: status code diferente de 200: " + strconv.Itoa(req.StatusCode))
	}
	defer req.Body.Close()
	res, err := io.ReadAll(req.Body)
	if err != nil {
		panic(err)
	}
	var data Localizacao
	err = json.Unmarshal(res, &data)
	if err != nil {
		panic(err)
	}
	return data
}

func buscaTemperatura(localizacao Localizacao) WeatherApiResponse {
	location := url.QueryEscape(localizacao.Localidade)
	address := "http://api.weatherapi.com/v1/current.json?key=e01c72f7886a4af1a7932746240704&q=" + location + "&aqi=no"
	log.Printf("URL WEATHER API: %s", address)
	req, err := http.Get(address)
	if err != nil {
		panic(err)
	}
	if req.StatusCode != http.StatusOK {
		panic("Erro ao fazer requisição para ViaCEP: status code diferente de 200: " + strconv.Itoa(req.StatusCode))
	}
	defer req.Body.Close()
	res, err := io.ReadAll(req.Body)
	if err != nil {
		panic(err)
	}
	var data WeatherApiResponse
	err = json.Unmarshal(res, &data)
	if err != nil {
		panic(err)
	}
	return data
}
