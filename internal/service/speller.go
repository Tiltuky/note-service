package service

import (
	"encoding/json"
	"net/http"
)

const apiURL = "https://speller.yandex.net/services/spellservice.json/checkText"

type SpellerClient struct{}

func NewSpellerClient() *SpellerClient {
	return &SpellerClient{}
}

func (c *SpellerClient) CheckText(text string) (string, error) {
	resp, err := http.PostForm(apiURL, map[string][]string{
		"text": {text},
	})
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	var result []struct {
		Code int      `json:"code"`
		Pos  int      `json:"pos"`
		Len  int      `json:"len"`
		Word string   `json:"word"`
		S    []string `json:"s"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", err
	}

	// Если нет ошибок, возвращаем исходный текст
	if len(result) == 0 {
		return text, nil
	}

	correctedText := []byte(text)

	// Исправляем слова на основе позиции и длины ошибки
	for i := len(result) - 1; i >= 0; i-- { // Применяем исправления с конца текста, чтобы избежать сдвигов
		correction := result[i]
		if len(correction.S) > 0 {
			correctedWord := correction.S[0]
			// Заменяем исходное слово исправленным
			correctedText = append(correctedText[:correction.Pos], append([]byte(correctedWord), correctedText[correction.Pos+correction.Len:]...)...)
		}
	}

	return string(correctedText), nil
}
