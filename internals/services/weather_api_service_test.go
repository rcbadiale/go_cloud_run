package services

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"testing"

	"regexp"

	"strings"

	"github.com/stretchr/testify/assert"
)

const (
	weatherAPIBody = `{ 
		"location": { 
			"name": "Florianópolis", 
			"region": "Santa Catarina", 
			"country": "Brasilien", 
			"lat": -27.58, 
			"lon": -48.57, 
			"tz_id": "America/Sao_Paulo", 
			"localtime_epoch": 1716578742, 
			"localtime": "2024-07-01 16:25" 
		}, 
		"current": {
			"last_updated_epoch": 1716578100, 
			"last_updated": "2024-07-01 16:15", 
			"temp_c": 21, 
			"temp_f": 69.8, 
			"is_day": 1, 
			"condition": { 
				"text": "Overcast", 
				"icon": "//cdn.weatherapi.com/weather/64x64/day/122.png",
				 "code": 1009 
			},
			"wind_mph": 8.1, 
			"wind_kph": 13,
			"wind_degree": 110, 
			"wind_dir": "ESE",
			"pressure_mb": 1008, 
			"pressure_in": 29.77, 
			"precip_mm": 0.04, 
			"precip_in": 0, 
			"humidity": 83, 
			"cloud": 100, 
			"feelslike_c": 21, 
			"feelslike_f": 69.8, 
			"vis_km": 10, 
			"vis_miles": 6,
			"uv": 5, 
			"gust_mph": 12.5, 
			"gust_kph": 20.2 
		} 
	}`
)

type mockWeatherApiHTTPClient struct{}

type brokenReader struct{}

func (br *brokenReader) Read(p []byte) (n int, err error) {
	return 0, fmt.Errorf("failed reading")
}

func (br *brokenReader) Close() error {
	return fmt.Errorf("failed closing")
}

func (m *mockWeatherApiHTTPClient) Get(url string) (*http.Response, error) {
	re := regexp.MustCompile(`^.*=.*=(.*)$`)
	match := re.FindStringSubmatch(url)
	switch match[1] {
	case "Florian%C3%B3polis":
		response := &http.Response{
			StatusCode: 200,
			Body:       io.NopCloser(bytes.NewBufferString(weatherAPIBody)),
			Header:     make(http.Header),
		}
		return response, nil
	case "Erroropolis":
		response := &http.Response{
			StatusCode: 400,
			Body:       io.NopCloser(bytes.NewBufferString("")),
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
	case "RequestFail":
		return nil, fmt.Errorf("error getting weather: 400")
	case "UnmarshalError":
		response := &http.Response{
			StatusCode: 200,
			Body:       io.NopCloser(bytes.NewBufferString(strings.Replace(weatherAPIBody, `"feelslike_f": 69.8`, `"feelslike_f": "quack"`, 1))),
			Header:     make(http.Header),
		}
		return response, nil
	default:
		return nil, nil
	}
}

func TestGetWeatherByCity(t *testing.T) {

	service := &WeatherAPIService{
		BaseHttpService: BaseHttpService{Client: &mockWeatherApiHTTPClient{}},
	}

	result, err := service.GetWeatherByCity("Florianópolis")

	// Assert there was no error
	assert.Nil(t, err)

	assert := assert.New(t)

	assert.Equal("Florianópolis", result.Location.Name)
	assert.Equal("Santa Catarina", result.Location.Region)
	assert.Equal("Brasilien", result.Location.Country)
	assert.Equal(float64(-27.58), result.Location.Lat)
	assert.Equal(float64(-48.57), result.Location.Lon)
	assert.Equal("America/Sao_Paulo", result.Location.TzID)
	assert.Equal(int(1716578742), result.Location.LocaltimeEpoch)
	assert.Equal("2024-07-01 16:25", result.Location.Localtime)

	assert.Equal(int(1716578100), result.Current.LastUpdatedEpoch)
	assert.Equal("2024-07-01 16:15", result.Current.LastUpdated)
	assert.Equal(float64(21), result.Current.TempC)
	assert.Equal(float64(69.8), result.Current.TempF)
	assert.Equal(1, result.Current.IsDay)
	assert.Equal("Overcast", result.Current.Condition.Text)
	assert.Equal("//cdn.weatherapi.com/weather/64x64/day/122.png", result.Current.Condition.Icon)
	assert.Equal(1009, result.Current.Condition.Code)
	assert.Equal(float64(8.1), result.Current.WindMph)
	assert.Equal(float64(13), result.Current.WindKph)
	assert.Equal(110, result.Current.WindDegree)
	assert.Equal("ESE", result.Current.WindDir)
	assert.Equal(float64(1008), result.Current.PressureMb)
	assert.Equal(float64(29.77), result.Current.PressureIn)
	assert.Equal(float64(0.04), result.Current.PrecipMm)
	assert.Equal(float64(0), result.Current.PrecipIn)
	assert.Equal(83, result.Current.Humidity)
	assert.Equal(100, result.Current.Cloud)
	assert.Equal(float64(21), result.Current.FeelslikeC)
	assert.Equal(float64(69.8), result.Current.FeelslikeF)
	assert.Equal(float64(10), result.Current.VisKm)
	assert.Equal(float64(6), result.Current.VisMiles)
	assert.Equal(float64(5), result.Current.Uv)
	assert.Equal(float64(12.5), result.Current.GustMph)
	assert.Equal(float64(20.2), result.Current.GustKph)
}

func TestGetWeatherNot200StatusCode(t *testing.T) {

	service := &WeatherAPIService{
		BaseHttpService: BaseHttpService{Client: &mockWeatherApiHTTPClient{}},
	}

	_, err := service.GetWeatherByCity("Erroropolis")

	assert := assert.New(t)
	assert.Equal(fmt.Errorf("error getting weather: 400"), err)

}

func TestGetWeatherBadBody(t *testing.T) {

	service := &WeatherAPIService{
		BaseHttpService: BaseHttpService{Client: &mockWeatherApiHTTPClient{}},
	}

	_, err := service.GetWeatherByCity("RequestFail")

	assert := assert.New(t)
	assert.Equal("error getting weather: 400", err.Error())
}

func TestGetWeatherUnmarshalError(t *testing.T) {

	service := &WeatherAPIService{
		BaseHttpService: BaseHttpService{Client: &mockWeatherApiHTTPClient{}},
	}

	_, err := service.GetWeatherByCity("UnmarshalError")

	assert := assert.New(t)
	assert.Equal("json: cannot unmarshal string into Go struct field WeatherAPIResponseCurrent.current.feelslike_f of type float64", err.Error())
}

func TestGetWeatherReaderError(t *testing.T) {

	service := &WeatherAPIService{
		BaseHttpService: BaseHttpService{Client: &mockWeatherApiHTTPClient{}},
	}

	_, err := service.GetWeatherByCity("BrokenReader")

	assert := assert.New(t)
	assert.Equal("failed reading", err.Error())
}
