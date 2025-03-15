package repository

import (
	"EduCommentSync/internal/models"
)

// GetColabLinksByWorkName Функция для получения ColabLinks по WorkName
func (r *repo) GetColabLinksByWorkName(workName string) ([]models.ColabLink, error) {
	var colabLinks []models.ColabLink

	result := r.dataBase.Where("work_name = ?", workName).Find(&colabLinks)
	if result.Error != nil {
		return nil, result.Error
	}

	return colabLinks, nil
}
