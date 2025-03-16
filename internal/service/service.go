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
