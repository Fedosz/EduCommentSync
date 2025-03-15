package repository

import (
	"EduCommentSync/internal/models"
)

func (r *repo) AddStudents(studentInfos []models.StudentInfo) ([]models.Student, error) {
	var students []models.Student

	for _, studentInfo := range studentInfos {
		student := models.Student{
			Name:     studentInfo.Name,
			MailHash: studentInfo.Mail,
		}

		result := r.dataBase.Create(&student)
		if result.Error != nil {
			return nil, result.Error
		}

		students = append(students, student)
	}

	return students, nil
}
