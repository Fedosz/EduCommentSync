package auth

import (
	"context"
	"fmt"
	"github.com/joho/godotenv"
	"log"
	"net/http"
	"os"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/drive/v3"
)

var oauth2Config *oauth2.Config

type OAuth2Token struct {
	*oauth2.Token
}

func init() {
	err := godotenv.Load() // Загружаем переменные из .env
	// Загружаем данные из credentials.json
	credentialsFile := os.Getenv("TOKEN_FILE")
	credentials, err := os.ReadFile(credentialsFile)
	if err != nil {
		log.Fatalf("Unable to read credentials file: %v", err)
	}

	// Создаем OAuth2 конфигурацию
	oauth2Config, err = google.ConfigFromJSON(credentials, drive.DriveScope)
	if err != nil {
		log.Fatalf("Unable to parse credentials: %v", err)
	}
}

// GetAuthURL генерирует URL для авторизации через Google OAuth2
func GetAuthURL() string {
	return oauth2Config.AuthCodeURL("", oauth2.AccessTypeOffline)
}

// GetClient получает клиента для работы с Google API, используя авторизационный код
func GetClient(code string) (*http.Client, error) {
	// Обмен кода на токен
	tok, err := oauth2Config.Exchange(context.Background(), code)
	if err != nil {
		return nil, fmt.Errorf("Unable to exchange code for token: %v", err)
	}

	// Создаем клиента
	client := oauth2Config.Client(context.Background(), tok)
	return client, nil
}
