package sheets

import (
	"EduCommentSync/internal/models"
	"fmt"
	"google.golang.org/api/sheets/v4"
	"net/http"
)

// GetSheetData retrieves data from a Google Sheet by spreadsheetId
func GetSheetData(client *http.Client, spreadsheetId string, tableName string) (*models.TableInfo, error) {
	service, err := sheets.New(client)
	if err != nil {
		return nil, fmt.Errorf("failed to create Sheets service: %v", err)
	}

	rangeData := "Sheet1!A1:Z1000"
	resp, err := service.Spreadsheets.Values.Get(spreadsheetId, rangeData).Do()
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve data from sheet: %v", err)
	}

	tableInfo := &models.TableInfo{
		Name:     tableName,
		Students: make([]models.StudentInfo, 0),
	}

	for i, row := range resp.Values {
		if i == 0 {
			continue
		}

		if len(row) >= 4 {
			info := models.StudentInfo{
				Name: row[1].(string),
				Mail: row[2].(string),
				Link: row[3].(string),
			}

			tableInfo.Students = append(tableInfo.Students, info)
		}
	}

	return tableInfo, nil
}
