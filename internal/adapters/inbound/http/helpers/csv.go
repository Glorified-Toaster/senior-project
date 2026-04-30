package helpers

import (
	"bytes"
	"encoding/csv"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"path/filepath"
	"slices"
	"strconv"
	"strings"
	"time"
	"uot-exam/internal/domain"
	"uot-exam/internal/ports"

	"github.com/gin-gonic/gin"
)

func MapQuestionCSVToStruct(file *multipart.FileHeader) ([]domain.Question, error) {

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

		if record[2] == string(domain.QuestionTypeImage) {
			continue
		}

		marks, err := strconv.ParseFloat(record[1], 64)
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

func ExportExamCSV(ctx *gin.Context, questions []domain.Question) {

	file := &bytes.Buffer{}
	csvWriter := csv.NewWriter(file)

	csvWriter.Write([]string{
		"Question Title",
		"Marks",
		"Question Type",
		"Question Text",
		"Choice 1",
		"Choice 2",
		"Choice 3",
		"Choice 4",
		"Correct Choice",
	})

	for _, question := range questions {
		if question.QuestionType == string(domain.QuestionTypeImage) {
			continue
		}
		correctChoice := ""
		choices := make([]string, 4)

		for i, choice := range question.Choices {
			if i >= 4 {
				break
			}
			choices[i] = choice.ChoiceText
			if choice.IsCorrect {
				correctChoice = choice.ChoiceText
			}
		}

		csvWriter.Write([]string{
			question.QuestionTitle,
			fmt.Sprintf("%.2f", question.Marks),
			question.QuestionType,
			question.QuestionText,
			choices[0],
			choices[1],
			choices[2],
			choices[3],
			correctChoice,
		})
	}

	csvWriter.Flush()
	if err := csvWriter.Error(); err != nil {
		ctx.Status(http.StatusInternalServerError)
		return
	}

	fileName := fmt.Sprintf("exam-%d.csv", time.Now().Unix())
	ctx.Header("Content-Description", "File Transfer")
	ctx.Header("Content-Disposition", fmt.Sprintf("attachment; filename=%s", fileName))
	ctx.Header("Content-Type", "text/csv")
	ctx.Data(http.StatusOK, "text/csv", file.Bytes())
}

func ParseUserCSV(file *multipart.FileHeader) ([]ports.CreateUserParams, error) {
	var users []ports.CreateUserParams

	fileReader, err := file.Open()
	if err != nil {
		return nil, err
	}
	defer fileReader.Close()

	csvReader := csv.NewReader(fileReader)
	csvReader.FieldsPerRecord = 4
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

		users = append(users, ports.CreateUserParams{
			Username: record[0],
			FullName: record[1],
			Role:     domain.UserRole(record[2]),
			Password: record[3],
		})
	}

	return users, nil
}

func MapStudentCSVToStruct(file *multipart.FileHeader) ([]string, error) {
	var users []string

	fileReader, err := file.Open()
	if err != nil {
		return nil, err
	}
	defer fileReader.Close()

	csvReader := csv.NewReader(fileReader)
	csvReader.FieldsPerRecord = 1
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

		users = append(users, record[0])
	}

	return users, nil
}

func ParseCSVFile(ctx *gin.Context, fileName string) (file *multipart.FileHeader, err error) {

	file, err = ctx.FormFile(fileName)

	if err != nil {
		return nil, err
	}

	if file == nil {
		return nil, &domain.ErrCSVUpload{Message: "File is required"}
	}

	if strings.ToLower(filepath.Ext(file.Filename)) != ".csv" {
		return nil, &domain.ErrCSVUpload{Message: "File must be a CSV file"}
	}

	if file.Size > 8<<20 {
		return nil, &domain.ErrCSVUpload{Message: "File size must be less than 8MB"}
	}

	return file, nil
}

func ExportStudentsCSV(ctx *gin.Context, students []domain.User) {
	file := &bytes.Buffer{}
	csvWriter := csv.NewWriter(file)

	csvWriter.Write([]string{
		"Full Name",
		"Username",
		"Role",
		"Status",
		"Created At",
	})

	for _, student := range students {
		status := "Disabled"
		if student.IsActive {
			status = "Enabled"
		}
		csvWriter.Write([]string{
			student.FullName,
			student.Username,
			string(student.Role),
			status,
			student.CreatedAt.Format("Jan 02, 2006"),
		})
	}

	csvWriter.Flush()
	if err := csvWriter.Error(); err != nil {
		ctx.Status(http.StatusInternalServerError)
		return
	}

	fileName := fmt.Sprintf("students-%d.csv", time.Now().Unix())
	ctx.Header("Content-Description", "File Transfer")
	ctx.Header("Content-Disposition", fmt.Sprintf("attachment; filename=%s", fileName))
	ctx.Header("Content-Type", "text/csv")
	ctx.Data(http.StatusOK, "text/csv", file.Bytes())
}
