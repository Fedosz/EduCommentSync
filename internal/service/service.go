package service

import (
	"EduCommentSync/internal/auth"
	"EduCommentSync/internal/config"
	"EduCommentSync/internal/models"
	"EduCommentSync/internal/processor"
	"EduCommentSync/internal/repository"
	"EduCommentSync/internal/sheets"
	"encoding/json"
	"fmt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
	"net/http"
	"strconv"
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
	http.HandleFunc("/loadXls", s.loadExcelFile)
	http.HandleFunc("/getExportsList", s.getExportList)
	http.HandleFunc("/getExport", s.getExportByID)
	http.HandleFunc("/addTeacher", s.addAuthor)

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
	data, err := sheets.GetSheetData(s.clientSync, fileID, fileName)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to get sheet data: %v", err), http.StatusInternalServerError)
		return
	}

	err = s.addInfo(data)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to add data: %v", err), http.StatusInternalServerError)
		return
	}

	links, err := s.repo.GetColabLinksByWorkName(data.Name)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to get colab links: %v", err), http.StatusBadRequest)
		return
	}

	err = s.processCommentsFromFiles(links)

	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to save comments: %v", err), http.StatusInternalServerError)
		return
	}

	err = s.process()

	if err != nil {
		http.Error(w, fmt.Sprintf("Error processing comments: %v", err), http.StatusInternalServerError)
		return
	}
}

func (s *Service) loadExcelFile(w http.ResponseWriter, r *http.Request) {
	s.clientMutex.Lock()
	defer s.clientMutex.Unlock()

	if s.clientSync == nil {
		http.Error(w, "Authentication required", http.StatusUnauthorized)
		return
	}

	comments, err := s.repo.GetComments()
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to get comments: %v", err), http.StatusInternalServerError)
		return
	}

	fullComments, err := s.EnrichComments(comments)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to enrich comments: %v", err), http.StatusInternalServerError)
		return
	}

	file := sheets.GenerateFile(fullComments)

	buf, err := file.WriteToBuffer()
	if err != nil {
		http.Error(w, "Не удалось создать файл", http.StatusInternalServerError)
		return
	}

	err = s.repo.AddExport(buf)
	if err != nil {
		http.Error(w, "Не удалось сохранить файл", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/octet-stream")
	w.Header().Set("Content-Disposition", "attachment; filename=Grades.xlsx")
	w.Header().Set("Content-Length", fmt.Sprint(buf.Len()))

	if _, err = w.Write(buf.Bytes()); err != nil {
		http.Error(w, "Не удалось отправить файл", http.StatusInternalServerError)
		return
	}
}

func (s *Service) addAuthor(w http.ResponseWriter, r *http.Request) {
	s.clientMutex.Lock()
	defer s.clientMutex.Unlock()

	if s.clientSync == nil {
		http.Error(w, "Authentication required", http.StatusUnauthorized)
		return
	}

	mail := r.URL.Query().Get("mail")
	if mail == "" {
		http.Error(w, "id is required", http.StatusBadRequest)
		return
	}

	err := s.repo.AddTeacher(mail)
	if err != nil {
		http.Error(w, "Error adding teacher mail", http.StatusInternalServerError)
		return
	}
}

func (s *Service) getExportList(w http.ResponseWriter, r *http.Request) {
	s.clientMutex.Lock()
	defer s.clientMutex.Unlock()

	if s.clientSync == nil {
		http.Error(w, "Authentication required", http.StatusUnauthorized)
		return
	}

	exports, err := s.repo.GetExports()
	if err != nil {
		http.Error(w, "Failed to get exports", http.StatusInternalServerError)
		return
	}

	var response []models.ExportResponse
	for _, export := range exports {
		response = append(response, models.ExportResponse{
			ID:         export.ID,
			ExportDate: export.ExportDate,
		})
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func (s *Service) getExportByID(w http.ResponseWriter, r *http.Request) {
	s.clientMutex.Lock()
	defer s.clientMutex.Unlock()

	if s.clientSync == nil {
		http.Error(w, "Authentication required", http.StatusUnauthorized)
		return
	}

	// Получаем данные из таблицы
	id := r.URL.Query().Get("id")
	if id == "" {
		http.Error(w, "id is required", http.StatusBadRequest)
		return
	}

	intID, err := strconv.Atoi(id)
	if err != nil {
		http.Error(w, "Wrong ID format", http.StatusBadRequest)
		return
	}

	export, err := s.repo.GetExportByID(int64(intID))
	if err != nil {
		http.Error(w, "Failed to get export", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/octet-stream")
	w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=Grades_%s.xlsx", id))

	if _, err = w.Write(export.FileData); err != nil {
		http.Error(w, "Не удалось отправить файл", http.StatusInternalServerError)
		return
	}
}
