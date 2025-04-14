package integration_test

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"io"
	"log"
	"net/http"
	"os"
	"testing"
	"time"
)

const (
	host           = "app:8080"
	healthPath     = "http://" + host + "/healthcheck"
	attempts       = 20
	requestTimeout = 5 * time.Second

	basePath = "http://" + host + "/api/v1"
)

var errHealthCheck = fmt.Errorf("url %s is not available", healthPath)

func doWebRequestWithTimeout(ctx context.Context, method, url string, body io.Reader) (*http.Response, error) {
	req, err := http.NewRequestWithContext(ctx, method, url, body)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")

	return http.DefaultClient.Do(req)
}

func doWebRequestWithTimeoutAuth(ctx context.Context, method, url, token string, body io.Reader) (*http.Response, error) {
	req, err := http.NewRequestWithContext(ctx, method, url, body)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)
	return http.DefaultClient.Do(req)
}

func getHealthCheck(url string) (int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), requestTimeout)

	defer cancel()

	resp, err := doWebRequestWithTimeout(ctx, http.MethodGet, url, http.NoBody)
	if err != nil {
		return -1, err
	}

	defer resp.Body.Close()

	return resp.StatusCode, nil
}

func healthCheck(attempts int) error {
	for attempts > 0 {
		statusCode, err := getHealthCheck(healthPath)
		if err != nil {
			return err
		}

		if statusCode == http.StatusOK {
			return nil
		}

		log.Printf("Integration tests: url %s is not available, attempts left: %d", healthPath, attempts)

		time.Sleep(time.Second)

		attempts--
	}

	return errHealthCheck
}

func getToken(ctx context.Context, role string) (string, error) {
	url := basePath + "/dummyLogin"
	payload := map[string]string{
		"role": role,
	}
	body, err := json.Marshal(payload)
	if err != nil {
		return "", fmt.Errorf("marshall body error: %v", err)
	}

	resp, err := doWebRequestWithTimeout(ctx, http.MethodPost, url, bytes.NewBuffer(body))
	if err != nil {
		return "", fmt.Errorf("failed to send request: %v", err)
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		b, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("bad status: %d, response: %s", resp.StatusCode, string(b))
	}

	var token string
	if err := json.NewDecoder(resp.Body).Decode(&token); err != nil {
		return "", fmt.Errorf("parse token error: %v", err)
	}
	return token, nil
}

func TestMain(m *testing.M) {
	err := healthCheck(attempts)
	if err != nil {
		log.Fatalf("Integration tests: host %s is not available: %s", host, err)
	}

	log.Printf("Integration tests: host %s is available", host)

	code := m.Run()
	os.Exit(code)
}

// TestIntegration реализует последовательный сценарий:
// 1. Получение токена модератора и создание ПВЗ.
// 2. Получение токена сотрудника ПВЗ и создание новой приемки.
// 3. Добавление 50 товаров к текущей приемке.
// 4. Закрытие приемки.
func TestIntegration(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	modToken, err := getToken(ctx, "moderator")
	if err != nil {
		t.Fatalf("Ошибка при получении токена модератора: %v", err)
	}

	pvzPayload := map[string]interface{}{
		"id":               uuid.New(),
		"registrationDate": time.Now(),
		"city":             "Москва",
	}
	pvzJSON, err := json.Marshal(pvzPayload)
	if err != nil {
		t.Fatalf("Не удалось получить данные ПВЗ: %v", err)
	}
	pvzURL := basePath + "/pvz"
	resp, err := doWebRequestWithTimeoutAuth(ctx, http.MethodPost, pvzURL, modToken, bytes.NewBuffer(pvzJSON))
	if err != nil {
		t.Fatalf("Ошибка выполнения запроса для создания ПВЗ: %v", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusCreated {
		b, _ := io.ReadAll(resp.Body)
		t.Fatalf("Ожидался статус 201 при создании ПВЗ, получен %d, ответ: %s", resp.StatusCode, string(b))
	}
	var createdPVZ map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&createdPVZ); err != nil {
		t.Fatalf("Ошибка получения ответа ПВЗ: %v", err)
	}
	pvzId, ok := createdPVZ["id"].(string)
	if !ok || pvzId == "" {
		t.Fatalf("В ответе отсутствует идентификатор ПВЗ")
	}

	empToken, err := getToken(ctx, "employee")
	if err != nil {
		t.Fatalf("Ошибка при получении токена сотрудника: %v", err)
	}

	receptionPayload := map[string]interface{}{
		"pvzId": pvzId,
	}
	receptionJSON, err := json.Marshal(receptionPayload)
	if err != nil {
		t.Fatalf("Ошибка получения данных приемки: %v", err)
	}
	receptionURL := basePath + "/receptions"
	resp, err = doWebRequestWithTimeoutAuth(ctx, http.MethodPost, receptionURL, empToken, bytes.NewBuffer(receptionJSON))
	if err != nil {
		t.Fatalf("Ошибка выполнения запроса для создания приемки: %v", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusCreated {
		b, _ := io.ReadAll(resp.Body)
		t.Fatalf("Ожидался статус 201 при создании приемки, получен %d, ответ: %s", resp.StatusCode, string(b))
	}

	productURL := basePath + "/products"
	for i := 1; i <= 50; i++ {
		productPayload := map[string]interface{}{
			"type":  "электроника",
			"pvzId": pvzId,
		}
		productJSON, err := json.Marshal(productPayload)
		if err != nil {
			t.Fatalf("Ошибка получения данных товара #%d: %v", i, err)
		}
		resp, err := doWebRequestWithTimeoutAuth(ctx, http.MethodPost, productURL, empToken, bytes.NewBuffer(productJSON))
		if err != nil {
			t.Fatalf("Ошибка выполнения запроса для товара #%d: %v", i, err)
		}
		if resp.StatusCode != http.StatusCreated {
			b, _ := io.ReadAll(resp.Body)
			resp.Body.Close()
			t.Fatalf("Ожидался статус 201 при добавлении товара #%d, получен %d, ответ: %s", i, resp.StatusCode, string(b))
		}
		resp.Body.Close()
	}

	closeReceptionURL := basePath + fmt.Sprintf("/pvz/%s/close_last_reception", pvzId)
	resp, err = doWebRequestWithTimeoutAuth(ctx, http.MethodPost, closeReceptionURL, empToken, nil)
	if err != nil {
		t.Fatalf("Ошибка выполнения запроса для закрытия приемки: %v", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		b, _ := io.ReadAll(resp.Body)
		t.Fatalf("Ожидался статус 200 при закрытии приемки, получен %d, ответ: %s", resp.StatusCode, string(b))
	}
}
