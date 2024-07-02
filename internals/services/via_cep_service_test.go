package services

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"net/http"
	"regexp"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

const CEPBody = `{
	"cep": "01001-000",
	"logradouro": "Praça da Sé",
	"complemento": "lado ímpar",
	"bairro": "Sé",
	"localidade": "São Paulo",
	"uf": "SP",
	"ibge": "3550308",
	"gia": "1004",
	"ddd": "11",
	"siafi": "7107"
}`

type mockViaCepHTTPClient struct{}

func (m *mockViaCepHTTPClient) Get(url string) (*http.Response, error) {
	re := regexp.MustCompile(`https://viacep.com.br/ws/(\w+)/json/`)
	match := re.FindStringSubmatch(url)
	log.Println("WOWOWO" + match[0])
	switch match[1] {
	case "01001000":
		response := &http.Response{
			StatusCode: 200,
			Body:       io.NopCloser(bytes.NewBufferString(CEPBody)),
			Header:     make(http.Header),
		}
		return response, nil
	case "Erroropolis":
		response := &http.Response{
			StatusCode: 400,
			Body:       io.NopCloser(bytes.NewBufferString(`{}`)),
			Header:     make(http.Header),
		}
		return response, nil
	case "BrokenReader":
		response := &http.Response{
			StatusCode: 200,
			Body:       &brokenReader{},
			Header:     make(http.Header),
		}
		return response, nil
	case "BadValue":
		response := &http.Response{
			StatusCode: 200,
			Body:       io.NopCloser(bytes.NewBufferString(`{"erro": "true"}`)),
			Header:     make(http.Header),
		}
		return response, nil
	case "RequestFail":
		return nil, fmt.Errorf("error getting weather: 400")
	case "UnmarshalError":
		response := &http.Response{
			StatusCode: 200,
			Body:       io.NopCloser(bytes.NewBufferString(strings.Replace(CEPBody, `"São Paulo"`, `000`, 1))),
			Header:     make(http.Header),
		}
		return response, nil
	default:
		return nil, nil
	}

}

func TestGetAddressByCEP(t *testing.T) {

	service := &ViaCEPService{
		BaseHttpService: BaseHttpService{Client: &mockViaCepHTTPClient{}},
	}

	response, err := service.GetAddressByCEP("01001000")

	// Assert there was no error
	assert.Nil(t, err)

	// Assert the response fields are correct
	assert.Equal(t, "01001-000", response.Cep)
	assert.Equal(t, "Praça da Sé", response.Logradouro)
	assert.Equal(t, "lado ímpar", response.Complemento)
	assert.Equal(t, "Sé", response.Bairro)
	assert.Equal(t, "São Paulo", response.Localidade)
	assert.Equal(t, "SP", response.Uf)
	assert.Equal(t, "3550308", response.Ibge)
	assert.Equal(t, "1004", response.Gia)
	assert.Equal(t, "11", response.Ddd)
	assert.Equal(t, "7107", response.Siafi)
}

func TestGetAddressReaderError(t *testing.T) {

	service := &ViaCEPService{
		BaseHttpService: BaseHttpService{Client: &mockViaCepHTTPClient{}},
	}

	_, err := service.GetAddressByCEP("BrokenReader")

	assert := assert.New(t)
	assert.Equal("failed reading", err.Error())
}

func TestGetAddressNot200StatusCode(t *testing.T) {

	service := &ViaCEPService{
		BaseHttpService: BaseHttpService{Client: &mockViaCepHTTPClient{}},
	}

	_, err := service.GetAddressByCEP("Erroropolis")

	assert := assert.New(t)
	assert.Equal(fmt.Errorf("invalid CEP provided"), err)

}

func TestGetAddressBadBody(t *testing.T) {

	service := &ViaCEPService{
		BaseHttpService: BaseHttpService{Client: &mockViaCepHTTPClient{}},
	}

	_, err := service.GetAddressByCEP("RequestFail")

	assert := assert.New(t)
	assert.Equal("error getting weather: 400", err.Error())
}

func TestGetAddressUnmarshalError(t *testing.T) {

	service := &ViaCEPService{
		BaseHttpService: BaseHttpService{Client: &mockViaCepHTTPClient{}},
	}

	resp, err := service.GetAddressByCEP("UnmarshalError")

	assert := assert.New(t)
	assert.Nil(resp)
	assert.Equal("invalid character '0' after object key:value pair", err.Error())
}

func TestGetAddressErrorinResponse(t *testing.T) {

	service := &ViaCEPService{
		BaseHttpService: BaseHttpService{Client: &mockViaCepHTTPClient{}},
	}

	resp, err := service.GetAddressByCEP("BadValue")

	assert := assert.New(t)
	assert.Nil(resp)
	assert.Equal("CEP not found", err.Error())
}
