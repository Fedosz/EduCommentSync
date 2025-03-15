package processor

import (
	"EduCommentSync/internal/models"
	"slices"
)

func ProcessComments(comments []models.RawComment, teachers []models.Teacher) []models.Comment {
	result := make([]models.Comment, 0)

	teachersNames := make([]string, 0, len(teachers))
	for _, t := range teachers {
		teachersNames = append(teachersNames, t.SysName)
	}

	for _, comment := range comments {
		if slices.Contains(teachersNames, comment.Author) {
			info := getCommentInfo(comment)
			if info == nil {
				continue
			}
			result = append(result, models.Comment{
				StudentID:  comment.StudentID,
				Text:       comment.Text,
				TaskNumber: info.TaskNumber,
				IsDone:     info.IsDone,
				WorkName:   comment.WorkName,
			})
		}
	}

	return result
}
