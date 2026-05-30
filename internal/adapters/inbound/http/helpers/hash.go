package helpers

import (
	"crypto/md5"
	"encoding/hex"
	"path/filepath"
	"strconv"
	"strings"
	"time"
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
	b.WriteString(q.QuestionText)
	b.WriteString(q.QuestionType)
	for _, c := range q.Choices {
		b.WriteString(c.ChoiceText)
	}
	return StringToMD5Hash([]byte(b.String()))
}

func HashFileName(filename string) string {
	extension := filepath.Ext(filename)
	name := strings.TrimSuffix(filename, extension)
	timestamp := time.Now().UnixNano()
	hashedName, _ := StringToMD5Hash([]byte(name + strconv.FormatInt(timestamp, 10)))
	return hashedName + extension
}
