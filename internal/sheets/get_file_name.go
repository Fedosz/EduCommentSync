package sheets

import (
	"context"
	"fmt"
	"google.golang.org/api/drive/v3"
	"google.golang.org/api/option"
	"net/http"
)

// GetFileName Функция для получения названия файла по fileId
func GetFileName(client *http.Client, fileId string) (string, error) {
	// Создаем контекст
	ctx := context.Background()

	// Создаем сервис Google Drive с использованием переданного клиента
	driveService, err := drive.NewService(ctx, option.WithHTTPClient(client))
	if err != nil {
		return "", fmt.Errorf("unable to create Drive client: %v", err)
	}

	// Получаем метаданные файла
	file, err := driveService.Files.Get(fileId).Do()
	if err != nil {
		return "", fmt.Errorf("unable to retrieve file metadata: %v", err)
	}

	// Возвращаем название файла
	return file.Name, nil
}
