# Temperatura por CEP com OTEL
Projeto composto por dois serviços:

- Serviço A - input validator

Recebe um CEP e faz a validação. Se tiver tudo OK, manda para o Serviço B.

- Serviço B - Orquestração

Busca a temperatura da localização do CEP.

## Execução servidor
`docker-compose up --build`

Porta padrão: 8080

## Api
Utilize a rota `POST /cep/` onde o valor do CEP deve ser somente números.

Ex: `curl -d '{"cep": "22460900"}' -H "Content-Type: application/json" -X POST http://localhost:8080/cep/`

Response ex: `{ "city: "São Paulo", "temp_C": 28.5, "temp_F": 28.5, "temp_K": 28.5 }`

### Status code
- 200: Success
- 422: Invalid zipcode
- 404: can not find zipcode
