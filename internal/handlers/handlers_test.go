package handlers

import (
	"fmt"
	"main/internal/config"
	"main/internal/models"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/google/uuid"
)

type MockDatabase struct {
	subscriptions []models.Subscribe
	errorMode     string
}

func (m *MockDatabase) Subscribe(sub models.Subscribe) (int, error) {
	if m.errorMode == "post" {
		return 0, fmt.Errorf("database error")
	}
	return 123, nil
}

func (m *MockDatabase) DeleteSubscribe(id uint) error {
	if m.errorMode == "delete" {
		return fmt.Errorf("database error")
	}
	return nil
}

func (m *MockDatabase) UpdateSubscribe(sub models.Subscribe) error {
	if m.errorMode == "put" {
		return fmt.Errorf("database error")
	}
	return nil
}

func (m *MockDatabase) GetSubscribe(f models.Filters) ([]models.Subscribe, error) {
	if m.errorMode == "get" {
		return nil, fmt.Errorf("database error")
	}
	return m.subscriptions, nil
}

func createTestHandlers(mockDB *MockDatabase) Handlers {
	return Handlers{
		Database: mockDB,
		Env:      config.Env{Port: "8080", Host: "localhost"},
	}
}

func TestSubscribeHandler(t *testing.T) {
	testUUID := uuid.New()
	futureDate := time.Now().AddDate(0, 2, 0).Format("01.2006")

	tests := []struct {
		name       string
		method     string
		body       string
		url        string
		errorMode  string
		wantStatus int
	}{
		// GET
		{"GET success", "GET", "", "/api/subscribe?service_name=Netflix", "", 200},
		{"GET empy params", "GET", "", "/api/subscribe", "", 400},
		{"GET invalid param", "GET", "", "/api/subscribe?qwe=123", "", 400},
		{"GET database error", "GET", "", "/api/subscribe", "get", 400},
		{"GET summary", "GET", "", "/api/subscribe?summary=true", "", 400},

		// POST
		{"POST success", "POST", fmt.Sprintf(`{"service_name":"Netflix","price":15,"user_id":"%s","start_date":"01.2025"}`, testUUID), "/api/subscribe", "", 201},
		{"POST invalid JSON", "POST", "invalid", "/api/subscribe", "", 400},
		{"POST invalid UUID", "POST", `{"service_name":"Netflix","price":15,"user_id":"invalid","start_date":"01.2025"}`, "/api/subscribe", "", 400},
		{"POST database error", "POST", fmt.Sprintf(`{"service_name":"Netflix","price":15,"user_id":"%s","start_date":"01.2025"}`, testUUID), "/api/subscribe", "post", 500},

		// PUT
		{"PUT success", "PUT", fmt.Sprintf(`{"id":1,"service_name":"Netflix","price":20,"user_id":"%s","start_date":"%s"}`, testUUID, futureDate), "/api/subscribe", "", 200},
		{"PUT invalid JSON", "PUT", "invalid", "/api/subscribe", "", 400},
		{"PUT invalid UUID", "PUT", `{"id":1,"service_name":"Netflix","price":20,"user_id":"invalid","start_date":"01.2025"}`, "/api/subscribe", "", 400},
		{"PUT database error", "PUT", fmt.Sprintf(`{"id":1,"service_name":"Netflix","price":20,"user_id":"%s","start_date":"%s"}`, testUUID, futureDate), "/api/subscribe", "put", 500},

		// DELETE
		{"DELETE success", "DELETE", `{"id":1}`, "/api/subscribe", "", 200},
		{"DELETE invalid ID", "DELETE", `{"id":-1}`, "/api/subscribe", "", 400},
		{"DELETE database error", "DELETE", `{"id":1}`, "/api/subscribe", "delete", 400},

		// Method not allowed
		{"PATCH not allowed", "PATCH", "", "/api/subscribe", "", 405},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockDB := &MockDatabase{
				subscriptions: []models.Subscribe{{Id: 1, ServiceName: "Netflix", Price: 15, UserId: testUUID}},
				errorMode:     tt.errorMode,
			}
			handlers := createTestHandlers(mockDB)

			body := strings.NewReader(tt.body)
			req := httptest.NewRequest(tt.method, tt.url, body)
			rr := httptest.NewRecorder()

			handlers.Subscribe(rr, req)

			if rr.Code != tt.wantStatus {
				t.Errorf("Expected status %d, got %d", tt.wantStatus, rr.Code)
			}

			if rr.Header().Get("Access-Control-Allow-Origin") != "*" {
				t.Error("CORS headers not set")
			}
		})
	}
}
