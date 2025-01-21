package server

import (
	"EduCommentSync/internal/auth"
	"EduCommentSync/internal/config"
	"EduCommentSync/internal/sheets"
	"fmt"
	"log"
	"net/http"
	"sync"
)

var (
	// Глобальная переменная для хранения токена
	clientSync  *http.Client
	clientMutex sync.Mutex
)

func StartServer() {
	cfg := config.LoadConfig()
	port := cfg.ServerPort

	http.HandleFunc("/auth", authHandler)                     // Маршрут для авторизации
	http.HandleFunc("/oauth2callback", oauth2CallbackHandler) // Обработчик ответа от Google
	http.HandleFunc("/getSheetData", getSheetDataHandler)     // Ручка для получения данных из Google Sheets

	// Запускаем сервер
	log.Printf("Starting server on port %s...", port)
	url := "http://localhost:8080/auth"
	fmt.Println("Перейдите по следующей ссылке для авторизации:", url)
	err := http.ListenAndServe(":"+port, nil)
	if err != nil {
		log.Fatal("Error starting server: ", err)
	}
}

func authHandler(w http.ResponseWriter, r *http.Request) {
	// Генерация URL авторизации
	authURL := auth.GetAuthURL()
	http.Redirect(w, r, authURL, http.StatusFound)
}

func oauth2CallbackHandler(w http.ResponseWriter, r *http.Request) {
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
	clientMutex.Lock()
	clientSync = client
	clientMutex.Unlock()

	// Отправляем сообщение, что авторизация успешна
	fmt.Fprintln(w, "Authentication successful! You can now use the application.")
}

func getSheetDataHandler(w http.ResponseWriter, r *http.Request) {
	// Проверяем, есть ли токен
	clientMutex.Lock()
	defer clientMutex.Unlock()

	if clientSync == nil {
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
	data, err := sheets.GetSheetData(clientSync, spreadsheetId)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to get sheet data: %v", err), http.StatusInternalServerError)
		return
	}

	// Отправляем данные в виде JSON
	w.Header().Set("Content-Type", "application/json")
	w.Write(data)
}
