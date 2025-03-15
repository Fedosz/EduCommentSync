package repository

import "EduCommentSync/internal/models"

func (r *repo) GetStudentById(id int) (*models.Student, error) {
	var student models.Student

	result := r.dataBase.Where("id = ?", id).Find(&student)
	if result.Error != nil {
		return nil, result.Error
	}

	return &student, nil
}
