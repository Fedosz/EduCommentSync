package models

import "time"

type StudentInfo struct {
	Name    string
	Surname string
	Mail    string
	Link    string
}

type TableInfo struct {
	Name     string
	Students []StudentInfo
}

// StudentComment - объединенная структура для студента и комментария
type StudentComment struct {
	CommentID  int
	StudentID  int
	Text       string
	TaskNumber int
	IsDone     bool
	WorkName   string
	Name       string
	SurName    string
	Mail       string
}

// ExportResponse Структура для ответа (ID и дата)
type ExportResponse struct {
	ID         int       `json:"id"`
	ExportDate time.Time `json:"export_date"`
}
