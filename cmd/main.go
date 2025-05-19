package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/dmmitrenko/weather-app/configs"
	"github.com/dmmitrenko/weather-app/internal/application"
	"github.com/dmmitrenko/weather-app/internal/infrastructure/cron"
	"github.com/dmmitrenko/weather-app/internal/infrastructure/emailing"
	weatherapi "github.com/dmmitrenko/weather-app/internal/infrastructure/weather-api"
	"github.com/dmmitrenko/weather-app/internal/repository"
	delivery "github.com/dmmitrenko/weather-app/internal/transport/http"
	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
)

func init() {
	log.SetOutput(os.Stdout)
	log.SetFlags(log.LstdFlags | log.Lshortfile)
}

func main() {
	cfg, err := configs.Load("configs/config.yaml")
	if err != nil {
		log.Fatal(err)
	}

	dsn := fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		cfg.DB.Host,
		cfg.DB.Port,
		cfg.DB.User,
		cfg.DB.Pass,
		cfg.DB.Name)

	db, err := sql.Open("postgres", dsn)
	if err != nil {
		log.Fatalf("db open: %v", err)
	}
	defer db.Close()
	if err := db.Ping(); err != nil {
		log.Fatalf("db ping: %v", err)
	}

	client := weatherapi.NewClient(cfg.WeatherAPI.Key)
	email_sender := emailing.NewSender(emailing.SmtpConfig{
		Host:     cfg.SMTP.Host,
		Port:     cfg.SMTP.Port,
		Username: cfg.SMTP.Username,
		Password: cfg.SMTP.Password,
		From:     cfg.SMTP.From,
	})

	subscriptionRepository := repository.NewSubscriptionRepository(db, cfg.Token.Secret)

	processor := &application.SubscriptionProcessor{Repo: subscriptionRepository, Sender: email_sender, WeatherProvider: client}
	subscriptionService := application.NewSubscriptionService(subscriptionRepository, email_sender)

	r := mux.NewRouter()
	delivery.NewWeatherHandler(r, client)
	delivery.NewSubscriptionHandler(r, subscriptionService)

	delivery.RegisterStatic(r)

	scheduler := cron.StartJobs(processor)
	defer scheduler.Stop()

	log.Printf("Listening on %s", cfg.Server.Address)
	log.Fatal(http.ListenAndServe(cfg.Server.Address, r))
}
