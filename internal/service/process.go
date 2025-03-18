package service

import (
	"EduCommentSync/internal/crypto"
	"EduCommentSync/internal/models"
	"EduCommentSync/internal/processor"
	"context"
	"fmt"
	"google.golang.org/api/drive/v3"
	"google.golang.org/api/option"
)

func (s *Service) process(workName string) error {
	rawComments, err := s.repo.GetRawComments(workName)
	if err != nil {
		return err
	}

	teachers, err := s.repo.GetTeachers()
	if err != nil {
		return err
	}

	for i, teacher := range teachers {
		teachers[i].SysName, err = crypto.Decrypt(s.cfg.SecretKey, teacher.SysName)
		if err != nil {
			return err
		}
	}
	comments := processor.ProcessComments(rawComments, teachers)

	err = s.repo.AddComments(comments)
	return err
}

func (s *Service) addInfo(info *models.TableInfo) error {
	students, err := s.repo.AddStudents(info.Students)
	if err != nil {
		return err
	}

	err = s.repo.AddColabLinks(info.Name, info.Students, students)
	return err
}

// processCommentsFromFiles получает комментарии из файлов на Google Colab
func (s *Service) processCommentsFromFiles(links []models.ColabLink) error {
	ctx := context.Background()
	srv, err := drive.NewService(ctx, option.WithHTTPClient(s.clientSync))
	if err != nil {
		return fmt.Errorf("не удалось создать сервис Google Drive: %v", err)
	}

	var rawComments []models.RawComment

	for _, link := range links {
		fileID, _ := processor.ExtractColabFileID(link.ColabLink)

		comments, err := srv.Comments.List(fileID).Fields("*").Do()
		if err != nil {
			return fmt.Errorf("не удалось получить комментарии: %v", err)
		}

		// Выводим комментарии
		for _, comment := range comments.Comments {
			rawComments = append(rawComments, models.RawComment{
				Text:      comment.Content,
				Author:    comment.Author.DisplayName,
				StudentID: link.StudentID,
				WorkName:  link.WorkName,
			})
		}
	}

	err = s.repo.AddRawComments(rawComments)

	return err
}

// EnrichComments объединяет комментарии и студентов
func (s *Service) EnrichComments(comments []models.Comment) ([]*models.StudentComment, error) {
	var combined []*models.StudentComment

	for _, comment := range comments {
		student, err := s.repo.GetStudentById(comment.StudentID)
		if err != nil {
			return nil, err
		}

		combined = append(combined, &models.StudentComment{
			CommentID:  comment.ID,
			StudentID:  comment.StudentID,
			Text:       comment.Text,
			TaskNumber: comment.TaskNumber,
			IsDone:     comment.IsDone,
			WorkName:   comment.WorkName,
			Name:       student.Name,
			SurName:    student.SurName,
			Mail:       student.Mail,
		})
	}

	err := s.DecryptStudents(combined)
	if err != nil {
		return nil, err
	}

	return combined, nil
}

func (s *Service) EncryptStudents(info *models.TableInfo) error {
	for i, data := range info.Students {
		encryptedMail, err := crypto.Encrypt(s.cfg.SecretKey, data.Mail)
		if err != nil {
			return err
		}
		encryptedName, err := crypto.Encrypt(s.cfg.SecretKey, data.Name)
		if err != nil {
			return err
		}
		encryptedSurname, err := crypto.Encrypt(s.cfg.SecretKey, data.Surname)
		if err != nil {
			return err
		}

		info.Students[i].Mail = encryptedMail
		info.Students[i].Name = encryptedName
		info.Students[i].Surname = encryptedSurname
	}

	return nil
}

func (s *Service) DecryptStudents(info []*models.StudentComment) error {
	for i, data := range info {
		decryptedMail, err := crypto.Decrypt(s.cfg.SecretKey, data.Mail)
		if err != nil {
			return err
		}
		decryptedName, err := crypto.Decrypt(s.cfg.SecretKey, data.Name)
		if err != nil {
			return err
		}
		decryptedSurname, err := crypto.Decrypt(s.cfg.SecretKey, data.SurName)
		if err != nil {
			return err
		}

		info[i].Mail = decryptedMail
		info[i].Name = decryptedName
		info[i].SurName = decryptedSurname

	}

	return nil
}
