package repository

import (
	"EduCommentSync/internal/models"
	"fmt"
)

func (r *repo) GetRawComments() ([]models.RawComment, error) {
	var rawComments []models.RawComment
	result := r.dataBase.Find(&rawComments)
	if result.Error != nil {
		return nil, fmt.Errorf("ошибка при получении данных: %v", result.Error)
	}

	return rawComments, nil
}
