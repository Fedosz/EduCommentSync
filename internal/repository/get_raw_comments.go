package repository

import (
	"EduCommentSync/internal/models"
	"fmt"
)

func (r *repo) GetRawComments(workName string) ([]models.RawComment, error) {
	var rawComments []models.RawComment
	result := r.dataBase.Where("work_name = ?", workName).Find(&rawComments)
	if result.Error != nil {
		return nil, fmt.Errorf("ошибка при получении данных: %v", result.Error)
	}

	return rawComments, nil
}
