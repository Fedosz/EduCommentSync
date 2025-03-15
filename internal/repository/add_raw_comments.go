package repository

import "EduCommentSync/internal/models"

func (r *repo) AddRawComments(comments []models.RawComment) error {
	tx := r.dataBase.Create(&comments)

	if tx.Error != nil {
		return tx.Error
	}
	return nil
}
