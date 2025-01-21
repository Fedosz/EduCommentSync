package sheets

import (
	"encoding/json"
	"fmt"
	"google.golang.org/api/sheets/v4"
	"net/http"
)

// GetSheetData retrieves data from a Google Sheet by spreadsheetId
func GetSheetData(client *http.Client, spreadsheetId string) ([]byte, error) {
	service, err := sheets.New(client)
	if err != nil {
		return nil, fmt.Errorf("failed to create Sheets service: %v", err)
	}

	// Получаем данные из диапазона (например, "Sheet1!A1:Z1000")
	rangeData := "Sheet1!A1:Z1000"
	resp, err := service.Spreadsheets.Values.Get(spreadsheetId, rangeData).Do()
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve data from sheet: %v", err)
	}

	// Преобразуем данные в JSON
	jsonData, err := json.Marshal(resp.Values)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal data to JSON: %v", err)
	}

	return jsonData, nil
}
