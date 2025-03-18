package processor

import (
	"EduCommentSync/internal/models"
	"sort"
)

// GetWorkNamesUnique возвращает список уникальных названий работ
func GetWorkNamesUnique(comments []*models.StudentComment) []string {
	workNames := make(map[string]struct{})
	for _, comment := range comments {
		workNames[comment.WorkName] = struct{}{}
	}

	uniqueWorkNames := make([]string, 0, len(workNames))
	for workName := range workNames {
		uniqueWorkNames = append(uniqueWorkNames, workName)
	}

	sort.Strings(uniqueWorkNames)
	return uniqueWorkNames
}
