package repository

import (
	"EduCommentSync/internal/models"
	"fmt"
)

func (r *repo) GetTeachers() ([]models.Teacher, error) {
	var teachers []models.Teacher

	result := r.dataBase.Find(&teachers)
	if result.Error != nil {
		return nil, fmt.Errorf("ошибка при получении данных: %v", result.Error)
	}

	return teachers, nil
}
