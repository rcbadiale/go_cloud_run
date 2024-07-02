package services

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
)

const (
	ViaCEP_URL = "https://viacep.com.br/ws/%s/json/"
)

type CEPService interface {
	GetAddressByCEP(cep string) (*ViaCEPResponse, error)
}

// ViaCEPService is a service to interact with the ViaCEP API
type ViaCEPService struct {
	BaseHttpService
}

type ViaCEPResponse struct {
	Cep         string `json:"cep"`
	Logradouro  string `json:"logradouro"`
	Complemento string `json:"complemento"`
	Bairro      string `json:"bairro"`
	Localidade  string `json:"localidade"`
	Uf          string `json:"uf"`
	Ibge        string `json:"ibge"`
	Gia         string `json:"gia"`
	Ddd         string `json:"ddd"`
	Siafi       string `json:"siafi"`
	Erro        string `json:"erro,omitempty"`
}

// NewViaCEPService creates a new ViaCEPService
func NewViaCEPService() CEPService {
	return &ViaCEPService{BaseHttpService{Client: &http.Client{}}}
}

// GetAddressByCEP returns the address for a given CEP
func (v *ViaCEPService) GetAddressByCEP(cep string) (*ViaCEPResponse, error) {

	resp, err := v.Client.Get(fmt.Sprintf(ViaCEP_URL, cep))
	if err != nil {
		log.Println("error getting address by CEP: ", err)
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Println("error reading response body: ", err)
		return nil, err
	} else if resp.StatusCode != 200 {
		return nil, ErrInvalidCEP
	}

	var viaCepResponse ViaCEPResponse
	err = json.Unmarshal(body, &viaCepResponse)
	if err != nil {
		log.Println("error on Unmarshal response body: ", err, string(body))
		return nil, err
	} else if viaCepResponse.Erro == "true" {
		log.Printf("error invalid address by CEP: %v\n", string(body))
		return nil, ErrCEPNotFound
	}

	return &viaCepResponse, nil
}
