package repository

import "EduCommentSync/internal/models"

func (r *repo) AddComments(comments []models.Comment) error {
	tx := r.dataBase.Create(&comments)

	if tx.Error != nil {
		return tx.Error
	}
	return nil
}
