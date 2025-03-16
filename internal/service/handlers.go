package service

import (
	"EduCommentSync/internal/processor"
	"EduCommentSync/internal/sheets"
	"fmt"
	"net/http"
)

func (s *Service) getSheetDataHandler(w http.ResponseWriter, r *http.Request) {
	s.clientMutex.Lock()
	defer s.clientMutex.Unlock()

	if s.clientSync == nil {
		http.Error(w, "Authentication required", http.StatusUnauthorized)
		return
	}

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

	err = s.process(data.Name)

	if err != nil {
		http.Error(w, fmt.Sprintf("Error processing comments: %v", err), http.StatusInternalServerError)
		return
	}
}

// addAuthor добавляет преподавателя
// @Summary Добавление преподавателя
// @Description Добавляет нового преподавателя по его email
// @Tags teachers
// @Accept json
// @Produce json
// @Param mail query string true "Email преподавателя"
// @Success 200 {string} string "Преподаватель успешно добавлен"
// @Failure 400 {string} string "mail is required"
// @Failure 401 {string} string "Authentication required"
// @Failure 500 {string} string "Error adding teacher mail"
// @Router /addTeacher [post]
func (s *Service) addAuthor(w http.ResponseWriter, r *http.Request) {
	s.clientMutex.Lock()
	defer s.clientMutex.Unlock()

	if s.clientSync == nil {
		http.Error(w, "Authentication required", http.StatusUnauthorized)
		return
	}

	mail := r.URL.Query().Get("mail")
	if mail == "" {
		http.Error(w, "mail is required", http.StatusBadRequest)
		return
	}

	err := s.repo.AddTeacher(mail)
	if err != nil {
		http.Error(w, "Error adding teacher mail", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Преподаватель успешно добавлен"))
}

func (s *Service) archive(w http.ResponseWriter, r *http.Request) {
	s.clientMutex.Lock()
	defer s.clientMutex.Unlock()

	if s.clientSync == nil {
		http.Error(w, "Authentication required", http.StatusUnauthorized)
		return
	}

	word := r.URL.Query().Get("approval")
	if word != "archive" {
		http.Error(w, "enter word 'archive' to continue", http.StatusBadRequest)
		return
	}

	err := s.repo.ArchiveData()
	if err != nil {
		http.Error(w, "Error adding teacher mail", http.StatusInternalServerError)
		return
	}
}
