package processor

import (
	"fmt"
	"strings"
)

// ExtractFileID Функция для извлечения fileId из ссылки
func ExtractFileID(url string) (string, error) {
	// Ищем начало fileId в ссылке
	start := strings.Index(url, "/d/")
	if start == -1 {
		return "", fmt.Errorf("invalid URL: fileId not found")
	}

	// Обрезаем начало ссылки до fileId
	url = url[start+3:]

	// Ищем конец fileId (либо "/", либо "?")
	end := strings.IndexAny(url, "/?")
	if end == -1 {
		return url, nil // Если нет "/" или "?", то вся оставшаяся строка — это fileId
	}

	// Возвращаем fileId
	return url[:end], nil
}
