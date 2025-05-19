package utils

import (
	"bytes"
	"fmt"
	"html/template"

	"github.com/dmmitrenko/weather-app/internal/domain"
)

var weatherTpl = template.Must(template.New("weatherUpdate").Parse(`
<!DOCTYPE html>
<html>
  <head>
    <meta charset="utf-8">
    <title>Weather Update</title>
  </head>
  <body>
    <p>Hello {{.Email}},</p>
    <p>Here is your weather update for <strong>{{.City}}</strong>:</p>
    <ul>
      <li>Temperature: {{printf "%.1f" .Temperature}} °C</li>
      <li>Humidity: {{.Humidity}}&#37;</li>
      <li>Description: {{.Description}}</li>
    </ul>
    <p>Have a great day!</p>
  </body>
</html>
`))

var confirmationTpl = template.Must(template.New("subscriptionConfirmation").Parse(`
<!DOCTYPE html>
<html>
  <head>
    <meta charset="utf-8">
    <title>Confirm Your Subscription</title>
  </head>
  <body>
    <p>Hello {{.Email}},</p>
    <p>Thank you for subscribing to our weather updates for <strong>{{.City}}</strong>.</p>
    <p>Please confirm your subscription using this token: {{.Token}}</p>
    <p>If you did not request this, you can safely ignore this email.</p>
    <p>Have a great day!</p>
  </body>
</html>
`))

type WeatherMessageData struct {
	Email       string
	City        string
	Temperature float64
	Humidity    int
	Description string
}

type ConfirmationMessageData struct {
	Email string
	City  string
	Token string
}

func BuildWeatherUpdateMessage(email string, city string, w domain.Weather) (subject, body string, err error) {
	subject = fmt.Sprintf("Weather update for %s", city)

	data := WeatherMessageData{
		Email:       email,
		City:        city,
		Temperature: w.Temperature,
		Humidity:    w.Humidity,
		Description: w.Description,
	}

	var buf bytes.Buffer
	if err := weatherTpl.Execute(&buf, data); err != nil {
		return "", "", fmt.Errorf("messagebuilder: template execution failed: %w", err)
	}
	return subject, buf.String(), nil
}

func BuildSubscriptionConfirmationMessage(email, city, token string) (subject, body string, err error) {
	subject = fmt.Sprintf("Please confirm your subscription for %s weather", city)

	data := ConfirmationMessageData{
		Email: email,
		City:  city,
		Token: token,
	}

	var buf bytes.Buffer
	if err := confirmationTpl.Execute(&buf, data); err != nil {
		return "", "", fmt.Errorf("messagebuilder: confirmation template execution failed: %w", err)
	}
	return subject, buf.String(), nil
}
