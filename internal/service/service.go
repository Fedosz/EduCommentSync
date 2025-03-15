package service

import (
	"EduCommentSync/internal/auth"
	"EduCommentSync/internal/config"
	"EduCommentSync/internal/processor"
	"EduCommentSync/internal/repository"
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
	repo        repository.Repository
	clientSync  *http.Client
	clientMutex sync.Mutex
}

func New() *Service {
	cfg := config.LoadConfig()

	return &Service{cfg: cfg}
}

func (s *Service) Run() error {
	db, err := gorm.Open(postgres.Open(s.cfg.DatabaseURL), &gorm.Config{})
	if err != nil {
		return err
	}
	s.repo, err = repository.NewRepository(db)
	if err != nil {
		return err
	}

	// Запуск сервера
	err = s.StartServer()
	if err != nil {
		return err
	}

	return nil
}

func (s *Service) StartServer() error {
	port := s.cfg.ServerPort

	http.HandleFunc("/auth", s.authHandler)                     // Маршрут для авторизации
	http.HandleFunc("/oauth2callback", s.oauth2CallbackHandler) // Обработчик ответа от Google
	http.HandleFunc("/getSheetData", s.getSheetDataHandler)     // Ручка для получения данных из Google Sheets

	// Запускаем сервер
	log.Printf("Starting server on port %s...", port)
	url := s.cfg.AuthURL
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
	link := r.URL.Query().Get("tables_link")
	if link == "" {
		http.Error(w, "tables_link is required", http.StatusBadRequest)
		return
	}

	fileID, err := processor.ExtractFileID(link)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	fileName, err := sheets.GetFileName(s.clientSync, fileID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Получаем данные из Google Sheets
	_, err = sheets.GetSheetData(s.clientSync, fileID, fileName)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to get sheet data: %v", err), http.StatusInternalServerError)
		return
	}

	// Отправляем данные в виде JSON
	w.Header().Set("Content-Type", "application/json")
	//w.Write(data)
}
