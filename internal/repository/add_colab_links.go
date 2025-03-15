package repository

import (
	"EduCommentSync/internal/models"
	"gorm.io/gorm"
)

func (r *repo) AddColabLinks(db *gorm.DB, workName string, studentInfos []models.StudentInfo, students []models.Student) error {
	for i, studentInfo := range studentInfos {
		colabLink := models.ColabLink{
			ColabLink: studentInfo.Link,
			WorkName:  workName,
			StudentID: students[i].ID,
		}

		result := db.Create(&colabLink)
		if result.Error != nil {
			return result.Error
		}
	}

	return nil
}
