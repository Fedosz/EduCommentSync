package service

import (
	"EduCommentSync/internal/auth"
	"EduCommentSync/internal/config"
	"EduCommentSync/internal/models"
	"EduCommentSync/internal/sheets"
	"fmt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
	"net/http"
	"sync"
)

type Service struct {
	cfg         config.Config
	dataBase    *gorm.DB
	clientSync  *http.Client
	clientMutex sync.Mutex
}

func New() *Service {
	cfg := config.LoadConfig()

	return &Service{cfg: cfg}
}

func (s *Service) Run() error {
	cfg := config.LoadConfig()
	db, err := gorm.Open(postgres.Open(cfg.DatabaseURL), &gorm.Config{})
	if err != nil {
		return err
	}
	s.dataBase = db
	models.AutoMigrate(s.dataBase)

	// Запуск сервера
	err = s.StartServer(cfg)
	if err != nil {
		return err
	}

	return nil
}

func (s *Service) StartServer(cfg config.Config) error {
	port := cfg.ServerPort

	http.HandleFunc("/auth", s.authHandler)                     // Маршрут для авторизации
	http.HandleFunc("/oauth2callback", s.oauth2CallbackHandler) // Обработчик ответа от Google
	http.HandleFunc("/getSheetData", s.getSheetDataHandler)     // Ручка для получения данных из Google Sheets

	// Запускаем сервер
	log.Printf("Starting server on port %s...", port)
	url := cfg.AuthURL
	fmt.Println("Перейдите по следующей ссылке для авторизации:", url)
	err := http.ListenAndServe(":"+port, nil)
	if err != nil {
		return err
	}
	return nil
}

func (s *Service) authHandler(w http.ResponseWriter, r *http.Request) {
	// Генерация URL авторизации
	authURL := auth.GetAuthURL()
	http.Redirect(w, r, authURL, http.StatusFound)
}

func (s *Service) oauth2CallbackHandler(w http.ResponseWriter, r *http.Request) {
	code := r.URL.Query().Get("code")
	if code == "" {
		http.Error(w, "Code not found", http.StatusBadRequest)
		return
	}

	// Получаем клиента Google API
	client, err := auth.GetClient(code)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to get client: %v", err), http.StatusInternalServerError)
		return
	}

	// Сохраняем токен в глобальную переменную для дальнейшего использования
	s.clientMutex.Lock()
	s.clientSync = client
	s.clientMutex.Unlock()

	// Отправляем сообщение, что авторизация успешна
	fmt.Fprintln(w, "Authentication successful! You can now use the application.")
}

func (s *Service) getSheetDataHandler(w http.ResponseWriter, r *http.Request) {
	// Проверяем, есть ли токен
	s.clientMutex.Lock()
	defer s.clientMutex.Unlock()

	if s.clientSync == nil {
		http.Error(w, "Authentication required", http.StatusUnauthorized)
		return
	}

	// Получаем данные из таблицы
	spreadsheetId := r.URL.Query().Get("spreadsheetId")
	if spreadsheetId == "" {
		http.Error(w, "SpreadsheetId is required", http.StatusBadRequest)
		return
	}

	// Получаем данные из Google Sheets
	data, err := sheets.GetSheetData(s.clientSync, spreadsheetId)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to get sheet data: %v", err), http.StatusInternalServerError)
		return
	}

	// Отправляем данные в виде JSON
	w.Header().Set("Content-Type", "application/json")
	w.Write(data)
}
