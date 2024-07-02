package client

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time-tracker/models"
)

func FetchUserData(user *models.User) error {
	// URL внешнего API
	apiUrl := "http://localhost:8080/api/userinfo"

	// Создание HTTP-клиента
	client := &http.Client{}

	// Формирование запроса
	req, err := http.NewRequest("GET", apiUrl, nil)
	if err != nil {
		return err
	}

	// Выполнение запроса
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Парсинг ответа
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected status: %s", resp.Status)
	}

	// Декодирование JSON-ответа в структуру User
	if err := json.NewDecoder(resp.Body).Decode(user); err != nil {
		return err
	}

	return nil
}
