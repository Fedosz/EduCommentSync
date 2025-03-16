package models

import (
	"gorm.io/gorm"
	"time"
)

// Student - таблица студентов
type Student struct {
	ID       int    `gorm:"primaryKey"`
	Name     string `gorm:"size:255;not null"`
	SurName  string `gorm:"size:255;not null"`
	MailHash string `gorm:"size:255;not null"`
}

// ColabLink - таблица ссылок на Colab
type ColabLink struct {
	ID        int    `gorm:"primaryKey"`
	ColabLink string `gorm:"size:255;not null"`
	StudentID int    `gorm:"not null"`
	WorkName  string `gorm:"not null"`
}

// RawComment - таблица сырых комментариев
type RawComment struct {
	ID        int    `gorm:"primaryKey"`
	Text      string `gorm:"type:text;not null"`
	Author    string
	StudentID int    `gorm:"not null"`
	WorkName  string `gorm:"not null"`
}

// Comment - таблица комментариев с привязкой к заданиям
type Comment struct {
	ID         int    `gorm:"primaryKey"`
	StudentID  int    `gorm:"not null"`
	Text       string `gorm:"type:text;not null"`
	TaskNumber int    `gorm:"not null"`
	IsDone     bool   `gorm:"not null"`
	WorkName   string `gorm:"not null"`
}

// Teacher - таблица комментариев с привязкой к заданиям
type Teacher struct {
	ID      int `gorm:"primaryKey"`
	SysName string
}

// CommentInfo - таблица комментариев с привязкой к заданиям
type CommentInfo struct {
	TaskNumber int
	IsDone     bool
}

// ExportFile Структура для хранения информации о выгрузке
type ExportFile struct {
	ID         int       `gorm:"primaryKey"`
	ExportDate time.Time `gorm:"not null"`
	FileData   []byte    `gorm:"type:bytea"`
}

// AutoMigrate выполняет миграции
func AutoMigrate(db *gorm.DB) error {
	err := db.AutoMigrate(&Student{}, &ColabLink{}, &RawComment{}, &Comment{}, &Teacher{}, &ExportFile{})
	if err != nil {
		return err
	}

	return nil
}
