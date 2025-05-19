package e2e_test

import (
	"os"
	"testing"
	"time"

	"github.com/gavv/httpexpect/v2"
)

func getBaseURL() string {
	if url := os.Getenv("E2E_BASE_URL"); url != "" {
		return url
	}
	return "http://app:8080"
}

func TestE2E(t *testing.T) {
	baseURL := getBaseURL()
	time.Sleep(5 * time.Second)

	e := httpexpect.New(t, baseURL)

	t.Run("GetWeather_MissingCity", func(t *testing.T) {
		e.GET("/api/weather").Expect().Status(400)
	})
	t.Run("GetWeather_UnknownCity", func(t *testing.T) {
		e.GET("/api/weather").
			WithQuery("city", "UnknownCityXYZ").
			Expect().
			Status(404)
	})
	t.Run("GetWeather_Success", func(t *testing.T) {
		obj := e.GET("/api/weather").
			WithQuery("city", "London").
			Expect().
			Status(200).
			JSON().Object()
		obj.ContainsKey("temperature")
		obj.ContainsKey("humidity")
		obj.ContainsKey("description")
	})

	payload := map[string]interface{}{"email": "test@example.com", "city": "London", "frequency": "hourly"}
	t.Run("Subscribe_BadRequest", func(t *testing.T) {
		e.POST("/api/subscribe").WithJSON(map[string]interface{}{}).Expect().Status(400)
	})
	t.Run("Subscribe_Success", func(t *testing.T) {
		e.POST("/api/subscribe").WithJSON(payload).Expect().Status(200)
	})
	t.Run("Subscribe_Conflict", func(t *testing.T) {
		e.POST("/api/subscribe").WithJSON(payload).Expect().Status(409)
	})

	t.Run("Confirm_InvalidToken", func(t *testing.T) {
		resp := e.GET("/api/confirm/{token}").
			WithPath("token", "").
			Expect()
		code := resp.Raw().StatusCode
		if code != 400 && code != 404 {
			t.Errorf("expected status 400 or 404, got %d", code)
		}
	})

	t.Run("Unsubscribe_InvalidToken", func(t *testing.T) {
		resp := e.GET("/api/unsubscribe/{token}").
			WithPath("token", "").
			Expect()
		code := resp.Raw().StatusCode
		if code != 400 && code != 404 {
			t.Errorf("expected status 400 or 404, got %d", code)
		}
	})
}
