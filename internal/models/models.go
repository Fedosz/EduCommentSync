package models

import "gorm.io/gorm"

// Students - таблица студентов
type Students struct {
	ID       int    `gorm:"primaryKey"`
	MailHash string `gorm:"size:255;not null"`
}

// ColabLinks - таблица ссылок на Colab
type ColabLinks struct {
	ID        int    `gorm:"primaryKey"`
	ColabLink string `gorm:"size:255;not null"`
	StudentID int    `gorm:"not null"`
	WorkID    int    `gorm:"not null"`
}

// RawComments - таблица сырых комментариев
type RawComments struct {
	ID        int    `gorm:"primaryKey"`
	Text      string `gorm:"type:text;not null"`
	Author    string
	StudentID int `gorm:"not null"`
	WorkID    int `gorm:"not null"`
}

// Comments - таблица комментариев с привязкой к заданиям
type Comments struct {
	ID         int `gorm:"primaryKey"`
	StudentID  int
	Text       string `gorm:"type:text;not null"`
	TaskNumber int    `gorm:"not null"`
	WorkID     int    `gorm:"not null"`
}

// Teachers - таблица комментариев с привязкой к заданиям
type Teachers struct {
	ID      int `gorm:"primaryKey"`
	SysName string
}

// AutoMigrate выполняет миграции
func AutoMigrate(db *gorm.DB) {
	err := db.AutoMigrate(&Students{}, &ColabLinks{}, &RawComments{}, &Comments{}, &Teachers{})
	if err != nil {
		return
	}
}
