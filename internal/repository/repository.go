package repository

import (
	"EduCommentSync/internal/models"
	"gorm.io/gorm"
)

type Repository interface {
	GetTeachers() ([]models.Teacher, error)
	AddTeacher(sysName string) error
	AddComments(comments []models.Comment) error
	GetRawComments() ([]models.RawComment, error)
	AddStudents(studentInfos []models.StudentInfo) ([]models.Student, error)
	AddColabLinks(workName string, studentInfos []models.StudentInfo, students []models.Student) error
	GetColabLinksByWorkName(workName string) ([]models.ColabLink, error)
	AddRawComments(comments []models.RawComment) error
	GetComments() ([]models.Comment, error)
	GetStudentById(id int) (*models.Student, error)
}

type repo struct {
	dataBase *gorm.DB
}

func NewRepository(db *gorm.DB) (*repo, error) {
	err := models.AutoMigrate(db)
	if err != nil {
		return nil, err
	}

	return &repo{dataBase: db}, nil
}
