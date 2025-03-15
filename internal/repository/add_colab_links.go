package repository

import (
	"EduCommentSync/internal/models"
)

func (r *repo) AddColabLinks(workName string, studentInfos []models.StudentInfo, students []models.Student) error {
	for i, studentInfo := range studentInfos {
		colabLink := models.ColabLink{
			ColabLink: studentInfo.Link,
			WorkName:  workName,
			StudentID: students[i].ID,
		}

		result := r.dataBase.Create(&colabLink)
		if result.Error != nil {
			return result.Error
		}
	}

	return nil
}
