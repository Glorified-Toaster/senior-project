package helpers

import (
	"encoding/csv"
	"io"
	"mime/multipart"
	"path/filepath"
	"slices"
	"strconv"
	"strings"
	"uot-exam/internal/domain"
	"uot-exam/web/templates/components/toast"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func MapCSVToStruct(file *multipart.FileHeader) ([]domain.Question, error) {

	var questions []domain.Question

	fileReader, err := file.Open()
	if err != nil {
		return nil, err
	}
	defer fileReader.Close()

	csvReader := csv.NewReader(fileReader)
	csvReader.FieldsPerRecord = 9
	csvReader.Read()

	for {
		record, err := csvReader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, err
		}
		if slices.Contains(record, "") {
			continue
		}
		marks, err := strconv.Atoi(record[1])
		if err != nil {
			return nil, err
		}
		questions = append(questions, domain.Question{
			QuestionTitle: record[0],
			Marks:         marks,
			QuestionType:  record[2],
			QuestionText:  record[3],
			Choices: []domain.Choice{
				{
					ChoiceText: record[4],
					IsCorrect:  "choice_1" == record[8],
				},
				{
					ChoiceText: record[5],
					IsCorrect:  "choice_2" == record[8],
				},
				{
					ChoiceText: record[6],
					IsCorrect:  "choice_3" == record[8],
				},
				{
					ChoiceText: record[7],
					IsCorrect:  "choice_4" == record[8],
				},
			},
		})
	}

	return questions, nil
}

func ParseCSVFile(ctx *gin.Context) (file *multipart.FileHeader, examID string, parsedUUID uuid.UUID, err error) {

	file, err = ctx.FormFile("file")

	if err != nil {
		ctx.Header("HX-Reswap", "none")
		Toast(ctx, "Upload Question CSV Failed", "Failed to get file", toast.VariantError)
		return
	}

	examID = ctx.Param("id")

	if IsTrimmedEmpty(examID) {
		return nil, "", uuid.Nil, &domain.ErrCSVUpload{Message: "Exam ID cannot be empty"}
	}

	parsedUUID, err = uuid.Parse(examID)

	if err != nil {
		return nil, "", uuid.Nil, &domain.ErrCSVUpload{Message: "Invalid exam ID"}
	}

	if strings.ToLower(filepath.Ext(file.Filename)) != ".csv" {
		return nil, "", uuid.Nil, &domain.ErrCSVUpload{Message: "File must be a CSV file"}
	}

	if file.Size > 8<<20 {
		return nil, "", uuid.Nil, &domain.ErrCSVUpload{Message: "File size must be less than 8MB"}
	}

	return file, examID, parsedUUID, nil
}
