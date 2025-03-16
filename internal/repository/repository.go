package repository

import (
	"EduCommentSync/internal/models"
	"bytes"
	"gorm.io/gorm"
	"log"
	"os"
)

type Repository interface {
	GetTeachers() ([]models.Teacher, error)
	AddTeacher(sysName string) error
	AddComments(comments []models.Comment) error
	GetRawComments(workName string) ([]models.RawComment, error)
	AddStudents(studentInfos []models.StudentInfo) ([]models.Student, error)
	AddColabLinks(workName string, studentInfos []models.StudentInfo, students []models.Student) error
	GetColabLinksByWorkName(workName string) ([]models.ColabLink, error)
	AddRawComments(comments []models.RawComment) error
	GetComments() ([]models.Comment, error)
	GetStudentById(id int) (*models.Student, error)
	AddExport(fileBytes *bytes.Buffer) error
	GetExports() ([]models.ExportFile, error)
	GetExportByID(id int64) (*models.ExportFile, error)
	ArchiveData() error
}

type repo struct {
	dataBase *gorm.DB
}

func NewRepository(db *gorm.DB) (*repo, error) {
	if _, err := os.Stat(".migration_done"); os.IsNotExist(err) {
		err = models.AutoMigrate(db)
		if err != nil {
			return nil, err
		}

		err = models.AutoMigrateArhive(db)
		if err != nil {
			return nil, err
		}

		file, err := os.Create(".migration_done")
		if err != nil {
			log.Fatalf("Failed to create migration flag: %v", err)
		}
		file.Close()

		log.Println("Migrations completed")
	} else {
		log.Println("Migrations already done")
	}

	return &repo{dataBase: db}, nil
}
