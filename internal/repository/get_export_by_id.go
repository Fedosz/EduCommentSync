package repository

import (
	"EduCommentSync/internal/models"
)

func (r *repo) GetExportByID(id int64) (*models.ExportFile, error) {
	var file models.ExportFile
	tx := r.dataBase.Where("id = ?", id).Find(&file)

	if tx.Error != nil {
		return nil, tx.Error
	}
	return &file, nil
}
