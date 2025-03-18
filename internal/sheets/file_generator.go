package sheets

import (
	"EduCommentSync/internal/models"
	"EduCommentSync/internal/processor"
	"fmt"
	"github.com/xuri/excelize/v2"
	"log"
	"sort"
	"strings"
)

func GenerateFile(comments []*models.StudentComment) *excelize.File {
	f := excelize.NewFile()
	defer func() {
		if err := f.Close(); err != nil {
			log.Fatal(err)
		}
	}()

	style, _ := f.NewStyle(&excelize.Style{
		Alignment: &excelize.Alignment{
			WrapText: true,
		},
	})

	uniqueWorkNames := processor.GetWorkNamesUnique(comments)

	headers := []string{"Фамилия", "Имя", "Адрес электронной почты"}
	for _, workName := range uniqueWorkNames {
		headers = append(headers, workName+" (Значение)", workName+" (Отзыв)")
	}

	for col, header := range headers {
		cell, _ := excelize.CoordinatesToCellName(col+1, 1)
		f.SetColWidth("Sheet1", "A", processor.GetColumnLetter(len(headers)), 30)
		f.SetCellValue("Sheet1", cell, header)
	}

	studentsMap := make(map[int][]*models.StudentComment)
	for _, comment := range comments {
		studentsMap[comment.StudentID] = append(studentsMap[comment.StudentID], comment)
	}

	row := 2
	for _, studentComments := range studentsMap {
		if len(studentComments) == 0 {
			continue
		}

		student := studentComments[0]
		f.SetCellValue("Sheet1", fmt.Sprintf("A%d", row), student.SurName)
		f.SetCellValue("Sheet1", fmt.Sprintf("B%d", row), student.Name)
		f.SetCellValue("Sheet1", fmt.Sprintf("C%d", row), student.Mail)

		sort.Slice(studentComments, func(i, j int) bool {
			return studentComments[i].TaskNumber < studentComments[j].TaskNumber
		})

		col := 4
		for _, workName := range uniqueWorkNames {
			var taskNumbers []int
			var feedbacks []string
			for _, comment := range studentComments {
				if comment.WorkName == workName {
					if comment.IsDone {
						taskNumbers = append(taskNumbers, comment.TaskNumber)
					}
					feedbacks = append(feedbacks, comment.Text)
				}
			}

			taskNumbersStr := strings.Join(strings.Fields(fmt.Sprint(taskNumbers)), " ")
			f.SetCellValue("Sheet1", fmt.Sprintf("%s%d", processor.GetColumnLetter(col), row), taskNumbersStr)

			feedbackText := strings.Join(feedbacks, "\n")
			cell := fmt.Sprintf("%s%d", processor.GetColumnLetter(col+1), row)
			f.SetCellValue("Sheet1", cell, feedbackText)
			f.SetCellStyle("Sheet1", cell, cell, style)

			col += 2
		}

		row++
	}

	return f
}
