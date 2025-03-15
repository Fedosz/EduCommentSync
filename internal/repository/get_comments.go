package repository

import (
	"EduCommentSync/internal/models"
	"fmt"
)

func (r *repo) GetComments() ([]models.Comment, error) {
	var Comments []models.Comment
	result := r.dataBase.Find(&Comments)
	if result.Error != nil {
		return nil, fmt.Errorf("ошибка при получении данных: %v", result.Error)
	}

	return Comments, nil
}
