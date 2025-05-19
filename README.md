# 🌦 Weather App
> A Go application for retrieving current weather and subscribing to weather updates via email.

### ⚙️ Requirements

- Docker + Docker Compose

### ▶️ Run Locally
To run the app locally:
```bash
docker compose up -d db
docker compose run --rm migrate
docker compose up -d app
```
- Web interface (subscription form): http://localhost:8080
- API base URL: http://localhost:8080/api

> ⚠️ **Note:**  
> The docker compose up command may sometimes not work as intended (migration may happen before the database is created)! Therefore, it is important to start containers sequentially, as specified in the instructions.

> ⚠️ **Note:**  
> The application requires a properly configured `.env` file containing credentials and API keys. Without it, email delivery and weather updates will not function correctly.

Example env file:
```env
SMTP_PASSWORD=""
WEATHERAPI_KEY="7"
SUBSCRIPTION_TOKEN_SECRET="It is not necessary to accept everything as true, one must only accept it as necessary"
```


### 📡 Integrations
####  Weather Data
This app fetches live weather information from **[WeatherAPI.com](https://www.weatherapi.com/)**. Make sure you sign up and obtain a free API key to use the service.


- You can get a key here: https://www.weatherapi.com/signup.aspx
- Set the key in your `.env` file:
          ```env
          WEATHERAPI_KEY=your_api_key_here
          ```

#### Email Delivery via Gmail SMTP
The application supports email delivery using Gmail SMTP.

Steps to enable Gmail SMTP:
1. Go to Google Account Security
2. Enable 2-Step Verification if not already enabled.
3. Under “Signing in to Google”, click App Passwords.
4. Create a new App Password for “Mail” + “Other (Custom name)” (e.g. Weather App)
5. Use the generated password in your .env:
        ```env
        SMTP_PASSWORD=generated_app_password
        ```
#### Security
**Subscription tokens are stored securely** in the database using HMAC-SHA256 hashing. To enable secure token hashing, you must define your own secret string in the `.env` file:
```env
SUBSCRIPTION_TOKEN_SECRET=your_custom_long_random_secret
```

### 🧑‍💼 What You Can Do

- **Check current weather** for any city using a simple web form or API request. ***(The application supports English only)**
- **Subscribe to weather updates** by entering your email, city, and preferred update frequency (`daily` or `hourly`). ***(Only one subscription per email)**
- **Receive email updates** with the latest weather forecasts for your selected city.
- **Unsubscribe anytime**

## 🗂 Project Structure
```
docker-compose.yml
Dockerfile
go.mod
go.sum
LICENSE
README.md
│
│
├───cmd
│       main.go
│
├───configs
│       config.go
│       config.yaml
│
├───docs
│       swagger.yaml
│
├───e2e
│       e2e_test.go
│
├───internal
│   ├───application
│   │       subscription_processor.go
│   │       subscription_service.go
│   │
│   ├───domain
│   │       errors.go
│   │       frequency.go
│   │       helpers.go
│   │       subscription.go
│   │       weather.go
│   │
│   ├───infrastructure
│   │   ├───cron
│   │   │       jobs.go
│   │   │
│   │   ├───emailing
│   │   │       email_sender.go
│   │   │
│   │   └───weather-api
│   │           client.go
│   │
│   ├───repository
│   │       subscription_repository.go
│   │
│   ├───transport
│   │   └───http
│   │       │   middleware.go
│   │       │   static.go
│   │       │   subscription_handler.go
│   │       │   weather_handler.go
│   │       │
│   │       └───static
│   │               index.html
│   │
│   └───utils
│           utils.go
│
└───scripts
        001_migrate_subscriptions.up.sql
```
