package handlers

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/rcbadiale/go-cloud-run/internals/services"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// Mock the services
type MockViaCEPService struct {
	mock.Mock
}

func (m *MockViaCEPService) GetAddressByCEP(cep string) (*services.ViaCEPResponse, error) {
	args := m.Called(cep)
	return args.Get(0).(*services.ViaCEPResponse), args.Error(1)
}

type MockWeatherAPIService struct {
	mock.Mock
}

func (m *MockWeatherAPIService) GetWeatherByCity(city string) (*services.WeatherAPIResponse, error) {
	args := m.Called(city)
	return args.Get(0).(*services.WeatherAPIResponse), args.Error(1)
}

func TestGetWeather(t *testing.T) {

	req, err := http.NewRequest("GET", "/weather/12345678", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	r := chi.NewRouter()
	mockViaCEPService := new(MockViaCEPService)
	mockViaCEPService.On(
		"GetAddressByCEP",
		"12345678").Return(
		&services.ViaCEPResponse{
			Localidade: "TestCity",
		}, nil,
	)

	mockWeatherService := new(MockWeatherAPIService)
	mockWeatherService.On(
		"GetWeatherByCity",
		"TestCity").Return(
		&services.WeatherAPIResponse{
			Current: services.WeatherAPIResponseCurrent{
				TempC: 10.0,
				TempF: 99.2,
			},
		},
		nil,
	)

	weatherHandler := WeatherHandler{
		CEPService:     mockViaCEPService,
		WeatherService: mockWeatherService,
	}
	r.Get("/weather/{zipCode}", weatherHandler.GetWeather)

	// Run the handler to get the response
	r.ServeHTTP(rr, req)

	// Check the status code is what we expect
	assert.Equal(t, http.StatusOK, rr.Code, "handler returned wrong status code")

	// Check the response body is what we expect
	expected := `{"temp_c":10,"temp_f":99.2,"temp_k":283.1}`
	assert.Equal(t, expected, strings.TrimRight(rr.Body.String(), "\n"), "handler returned unexpected body")
	assert.Equal(t, 200, rr.Result().StatusCode, "handler returned unexpected statusCode")
}
