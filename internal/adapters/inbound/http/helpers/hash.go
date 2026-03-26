package helpers

import (
	"crypto/md5"
	"encoding/hex"
	"strconv"
	"strings"
	"uot-exam/internal/domain"
)

func StringToMD5Hash(text []byte) (string, error) {
	hasher := md5.New()
	hasher.Write(text)
	return hex.EncodeToString(hasher.Sum(nil)), nil
}

func BuildQuestionChecksum(examID string, q domain.Question) (string, error) {
	var b strings.Builder
	b.WriteString(examID)
	b.WriteString(q.QuestionTitle)
	b.WriteString(q.QuestionText)
	b.WriteString(q.QuestionType)
	b.WriteString(strconv.Itoa(q.Marks))
	for _, c := range q.Choices {
		b.WriteString(c.ChoiceText)
	}
	return StringToMD5Hash([]byte(b.String()))
}
