package weatherapi

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"

	"github.com/dmmitrenko/weather-app/internal/domain"
)

const endpoint = "https://api.weatherapi.com/v1/current.json"

type apiError struct {
	Error struct {
		Code    int    `json:"code"`
		Message string `json:"message"`
	} `json:"error"`
}

type Client struct {
	apiKey     string
	httpClient *http.Client
}

func NewClient(apiKey string) *Client {
	return &Client{
		apiKey: apiKey,
		httpClient: &http.Client{
			Timeout: 5 * time.Second,
		},
	}
}

type apiResponse struct {
	Location struct {
		Name      string `json:"name"`
		Region    string `json:"region"`
		Country   string `json:"country"`
		Localtime string `json:"localtime"`
	} `json:"location"`
	Current struct {
		LastUpdated string  `json:"last_updated"`
		TempC       float64 `json:"temp_c"`
		Humidity    int     `json:"humidity"`
		Condition   struct {
			Text string `json:"text"`
			Icon string `json:"icon"`
			Code int    `json:"code"`
		} `json:"condition"`
	} `json:"current"`
}

func (c *Client) GetCurrentWeather(ctx context.Context, city string) (domain.Weather, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, endpoint, nil)

	if err != nil {
		return domain.Weather{}, err
	}

	q := req.URL.Query()
	q.Set("key", c.apiKey)
	q.Set("q", city)
	q.Set("aqi", "yes")
	req.URL.RawQuery = q.Encode()

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return domain.Weather{}, err
	}
	defer resp.Body.Close()

	raw, err := io.ReadAll(resp.Body)
	if err != nil {
		return domain.Weather{}, err
	}

	var ferr apiError
	if err := json.Unmarshal(raw, &ferr); err == nil && ferr.Error.Message != "" {
		log.Printf("WeatherAPI error %d: %s", ferr.Error.Code, ferr.Error.Message)
		if ferr.Error.Code == 1006 {
			return domain.Weather{}, domain.ErrCityNotFound
		}
		return domain.Weather{}, fmt.Errorf("weatherapi error %d: %s", ferr.Error.Code, ferr.Error.Message)
	}

	resp.Body = io.NopCloser(bytes.NewBuffer(raw))

	switch resp.StatusCode {
	case http.StatusOK:
	case http.StatusBadRequest:
		return domain.Weather{}, fmt.Errorf("bad request to WeatherAPI")
	default:
		return domain.Weather{}, fmt.Errorf("weatherapi: unexpected status %d", resp.StatusCode)
	}

	var apiRes apiResponse
	if err := json.NewDecoder(resp.Body).Decode(&apiRes); err != nil {
		return domain.Weather{}, err
	}

	w := domain.Weather{
		City:        apiRes.Location.Name,
		Region:      apiRes.Location.Region,
		Country:     apiRes.Location.Country,
		LocalTime:   parseTime(apiRes.Location.Localtime),
		LastUpdated: parseTime(apiRes.Current.LastUpdated),
		Temperature: apiRes.Current.TempC,
		Humidity:    apiRes.Current.Humidity,
		Description: apiRes.Current.Condition.Text,
	}
	return w, nil
}

func parseTime(ts string) time.Time {
	t, _ := time.Parse("2006-01-02 15:04", ts)
	return t
}
