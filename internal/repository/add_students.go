package repository

import (
	"EduCommentSync/internal/models"
	"errors"
	"gorm.io/gorm"
)

func (r *repo) AddStudents(studentInfos []models.StudentInfo) ([]models.Student, error) {
	var students []models.Student

	// Проходим по каждому студенту
	for _, studentInfo := range studentInfos {
		var student models.Student

		// Проверяем, существует ли студент в базе данных
		result := r.dataBase.Where("mail = ?", studentInfo.Mail).First(&student)
		if result.Error != nil {
			if errors.Is(result.Error, gorm.ErrRecordNotFound) {
				student = models.Student{
					Name:    studentInfo.Name,
					SurName: studentInfo.Surname,
					Mail:    studentInfo.Mail,
				}

				result = r.dataBase.Create(&student)
				if result.Error != nil {
					return nil, result.Error
				}
			} else {
				return nil, result.Error
			}
		}

		students = append(students, student)
	}

	return students, nil
}
