package processor

import (
	"EduCommentSync/internal/models"
	"fmt"
	"regexp"
	"strconv"
)

func getCommentInfo(comment models.RawComment) *models.CommentInfo {
	re := regexp.MustCompile(`^(\d{1,2})\+?`)

	// Поиск совпадений
	matches := re.FindStringSubmatch(comment.Text)
	if matches == nil {
		fmt.Println("Строка не соответствует формату.")
		return nil
	}

	// Извлечение номера задания
	taskNumber, err := strconv.Atoi(matches[1])
	if err != nil {
		fmt.Println("Ошибка при преобразовании номера задания:", err)
		return nil
	}

	// Проверка наличия знака "+"
	hasPlus := len(matches[0]) > len(matches[1])

	return &models.CommentInfo{
		TaskNumber: taskNumber,
		IsDone:     hasPlus,
	}
}
