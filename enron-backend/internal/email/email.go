package email

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"
)

type Email struct {
	Id        string `json:"id"`
	From      string `json:"from"`
	To        string `json:"to"`
	Subject   string `json:"subject"`
	Body      string `json:"body"`
	Date      string `json:"date"`
	Highlight struct {
		Body    []string `json:"body"`
		Subject []string `json:"subject"`
		From    []string `json:"from"`
		To      []string `json:"to"`
	} `json:"highlight"`
}

func ReadEmail(filePath string) (Email, error) {
	content, err := os.ReadFile(filePath)
	if err != nil {
		return Email{}, fmt.Errorf("no se pudo leer el archivo %s: %v", filePath, err)
	}

	parts := strings.SplitN(string(content), "\r\n\r\n", 2)
	if len(parts) < 2 {
		return Email{}, fmt.Errorf("formato de email inválido en el archivo %s", filePath)
	}
	headers := strings.Split(strings.ReplaceAll(parts[0], "\n\t\r", ""), "\r\n")
	body := parts[1]
	email := Email{}
	for _, line := range headers {
		switch {
		case strings.HasPrefix(line, "From:"):
			email.From = strings.TrimSpace(strings.TrimPrefix(line, "From:"))

		case strings.HasPrefix(line, "To:"):
			email.To = strings.TrimSpace(strings.TrimPrefix(line, "To:"))

		case strings.HasPrefix(line, "Subject:"):
			email.Subject = strings.TrimSpace(strings.TrimPrefix(line, "Subject:"))

		case strings.HasPrefix(line, "Date:"):
			dateValue := strings.TrimSpace(strings.TrimPrefix(line, "Date:"))
			if idx := strings.IndexByte(dateValue, '('); idx != -1 {
				dateValue = dateValue[:idx]
			}
			dateValue = strings.TrimSpace(dateValue)

			t, err := time.Parse("Mon, 2 Jan 2006 15:04:05 -0700", dateValue)
			if err != nil {
				return Email{}, fmt.Errorf("error parsing date: %v", err)
			}
			email.Date = t.Format("Mon, 02 Jan 2006 15:04:05 -0700")
		}
	}
	email.Body = body
	return email, nil
}

func ReadEmails(dirPath string) ([]Email, error) {
	var emails []Email
	err := filepath.Walk(dirPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() {
			email, err := ReadEmail(path)
			if err != nil {
				return nil
			}
			emails = append(emails, email)
		}
		return nil
	})
	return emails, err
}
