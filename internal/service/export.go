package service

import (
	"EduCommentSync/internal/models"
	"EduCommentSync/internal/sheets"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
)

// loadExcelFile генерирует и возвращает Excel-файл с комментариями
// @Summary Генерация Excel-файла
// @Description Генерирует Excel-файл на основе комментариев и возвращает его для скачивания
// @Tags files
// @Accept json
// @Produce application/octet-stream
// @Success 200 {file} file "Excel-файл с комментариями"
// @Failure 401 {string} string "Authentication required"
// @Failure 500 {string} string "Failed to get comments"
// @Failure 500 {string} string "Failed to enrich comments"
// @Failure 500 {string} string "Не удалось создать файл"
// @Failure 500 {string} string "Не удалось сохранить файл"
// @Failure 500 {string} string "Не удалось отправить файл"
// @Router /loadXls [get]
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

// getExportList возвращает список экспортов
// @Summary Получение списка экспортов
// @Description Возвращает список всех экспортов из базы данных
// @Tags exports
// @Accept json
// @Produce json
// @Success 200 {array} models.ExportResponse "Список экспортов"
// @Failure 401 {string} string "Authentication required"
// @Failure 500 {string} string "Failed to get exports"
// @Router /getExportsList [get]
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

// getExportByID возвращает экспорт по ID
// @Summary Получение экспорта по ID
// @Description Возвращает экспорт по указанному ID
// @Tags exports
// @Accept json
// @Produce json
// @Param id query int true "ID экспорта"
// @Success 200 "Файл успешно скачан"
// @Router /getExport [get]
func (s *Service) getExportByID(w http.ResponseWriter, r *http.Request) {
	s.clientMutex.Lock()
	defer s.clientMutex.Unlock()

	if s.clientSync == nil {
		http.Error(w, "Authentication required", http.StatusUnauthorized)
		return
	}

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
