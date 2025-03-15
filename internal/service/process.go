package service

import (
	"EduCommentSync/internal/models"
	"EduCommentSync/internal/processor"
	"context"
	"fmt"
	"google.golang.org/api/drive/v3"
	"google.golang.org/api/option"
)

func (s *Service) process() error {
	rawComments, err := s.repo.GetRawComments()
	if err != nil {
		return err
	}

	teachers, err := s.repo.GetTeachers()
	if err != nil {
		return err
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
				Author:    comment.Author.EmailAddress,
				StudentID: link.StudentID,
				WorkName:  link.WorkName,
			})
		}
	}

	err = s.repo.AddRawComments(rawComments)

	return err
}

// EnrichComments объединяет комментарии и студентов
func (s *Service) EnrichComments(comments []models.Comment) ([]models.StudentComment, error) {
	var combined []models.StudentComment

	for _, comment := range comments {
		student, err := s.repo.GetStudentById(comment.StudentID)
		if err != nil {
			return nil, err
		}

		combined = append(combined, models.StudentComment{
			CommentID:  comment.ID,
			StudentID:  comment.StudentID,
			Text:       comment.Text,
			TaskNumber: comment.TaskNumber,
			IsDone:     comment.IsDone,
			WorkName:   comment.WorkName,
			Name:       student.Name,
			SurName:    student.SurName,
			MailHash:   student.MailHash,
		})
	}

	return combined, nil
}
