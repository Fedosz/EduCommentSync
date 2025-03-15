package repository

import "EduCommentSync/internal/models"

func (r *repo) AddTeacher(sysName string) error {
	tx := r.dataBase.Create(&models.Teacher{
		SysName: sysName,
	})

	if tx.Error != nil {
		return tx.Error
	}
	return nil
}
