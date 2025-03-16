package repository

import (
	"EduCommentSync/internal/models"
	"gorm.io/gorm"
	"log"
)

func (r *repo) ArchiveData() error {
	tx := r.dataBase.Begin()
	if tx.Error != nil {
		return tx.Error
	}

	defer func() {
		if rec := recover(); rec != nil {
			tx.Rollback()
			log.Println("Transaction rolled back due to panic:", rec)
		}
	}()

	err := r.archiveStudents(tx)
	if err != nil {
		tx.Rollback()
		return err
	}

	err = r.archiveColabLinks(tx)
	if err != nil {
		tx.Rollback()
		return err
	}

	err = r.archiveRawComments(tx)
	if err != nil {
		tx.Rollback()
		return err
	}

	err = r.archiveComments(tx)
	if err != nil {
		tx.Rollback()
		return err
	}

	if err = tx.Commit().Error; err != nil {
		log.Println("Failed to commit transaction:", err)
		return err
	}
	return nil
}

func (r *repo) archiveStudents(tx *gorm.DB) error {
	var students []models.Student
	result := r.dataBase.Find(&students)
	if result.Error != nil {
		tx.Rollback()
		return result.Error
	}

	for _, student := range students {
		archive := models.StudentArchive{
			Name:     student.Name,
			SurName:  student.SurName,
			MailHash: student.MailHash,
		}
		if err := tx.Create(&archive).Error; err != nil {
			tx.Rollback()
			return err
		}
	}

	if err := tx.Exec("DELETE FROM students").Error; err != nil {
		tx.Rollback()
		return err
	}

	return nil
}

func (r *repo) archiveColabLinks(tx *gorm.DB) error {
	var colabLinks []models.ColabLink
	result := r.dataBase.Find(&colabLinks)
	if result.Error != nil {
		tx.Rollback()
		return result.Error
	}

	for _, link := range colabLinks {
		archive := models.ColabLinkArchive{
			ColabLink: link.ColabLink,
			StudentID: link.StudentID,
			WorkName:  link.WorkName,
		}
		if err := tx.Create(&archive).Error; err != nil {
			tx.Rollback()
			return err
		}
	}

	if err := tx.Exec("DELETE FROM colab_links").Error; err != nil {
		tx.Rollback()
		return err
	}

	return nil
}

func (r *repo) archiveRawComments(tx *gorm.DB) error {
	var rawComments []models.RawComment
	result := r.dataBase.Find(&rawComments)
	if result.Error != nil {
		tx.Rollback()
		return result.Error
	}

	for _, comment := range rawComments {
		archive := models.RawCommentArchive{
			Text:      comment.Text,
			Author:    comment.Author,
			StudentID: comment.StudentID,
			WorkName:  comment.WorkName,
		}
		if err := tx.Create(&archive).Error; err != nil {
			tx.Rollback()
			return err
		}
	}

	if err := tx.Exec("DELETE FROM raw_comments").Error; err != nil {
		tx.Rollback()
		return err
	}

	return nil
}

func (r *repo) archiveComments(tx *gorm.DB) error {
	var comments []models.Comment
	result := r.dataBase.Find(&comments)
	if result.Error != nil {
		tx.Rollback()
		return result.Error
	}

	for _, comment := range comments {
		archive := models.CommentArchive{
			StudentID:  comment.StudentID,
			Text:       comment.Text,
			TaskNumber: comment.TaskNumber,
			IsDone:     comment.IsDone,
			WorkName:   comment.WorkName,
		}
		if err := tx.Create(&archive).Error; err != nil {
			tx.Rollback()
			return err
		}
	}

	if err := tx.Exec("DELETE FROM comments").Error; err != nil {
		tx.Rollback()
		return err
	}

	return nil
}
