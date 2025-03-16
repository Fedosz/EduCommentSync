package service

import (
	_ "EduCommentSync/docs"
	"EduCommentSync/internal/config"
	"EduCommentSync/internal/repository"
	"fmt"
	"github.com/swaggo/http-swagger"
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

// @title EduCommentSync API
// @version 1.0
// @description API для синхронизации комментариев и работы с Google Sheets.
// @termsOfService http://swagger.io/terms/
// @contact.name API Support
// @contact.email support@educommentsync.com
// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html
// @host localhost:8080
// @BasePath /
func (s *Service) StartServer() error {
	port := s.cfg.ServerPort

	http.HandleFunc("/swagger/", httpSwagger.WrapHandler)

	http.HandleFunc("/auth", s.authHandler)
	http.HandleFunc("/oauth2callback", s.oauth2CallbackHandler)
	http.HandleFunc("/getSheetData", s.getSheetDataHandler)
	http.HandleFunc("/loadXls", s.loadExcelFile)
	http.HandleFunc("/getExportsList", s.getExportList)
	http.HandleFunc("/getExport", s.getExportByID)
	http.HandleFunc("/addTeacher", s.addAuthor)
	http.HandleFunc("/archive", s.archive)

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
