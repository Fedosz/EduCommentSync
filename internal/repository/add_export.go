package repository

import (
	"EduCommentSync/internal/models"
	"bytes"
	"time"
)

func (r *repo) AddExport(fileBytes *bytes.Buffer) error {
	file := models.ExportFile{
		ExportDate: time.Now(),
		FileData:   fileBytes.Bytes(),
	}
	tx := r.dataBase.Create(&file)

	if tx.Error != nil {
		return tx.Error
	}
	return nil
}
