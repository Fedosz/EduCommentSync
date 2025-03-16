package repository

import (
	"EduCommentSync/internal/models"
)

func (r *repo) GetExports() ([]models.ExportFile, error) {
	var files []models.ExportFile
	tx := r.dataBase.Find(&files)

	if tx.Error != nil {
		return nil, tx.Error
	}
	return files, nil
}
