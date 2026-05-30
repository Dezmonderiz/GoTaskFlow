package handler_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"GoTaskFlow/internal/handler"
	"GoTaskFlow/internal/model"
	"GoTaskFlow/internal/repository"
	"GoTaskFlow/internal/service"

	"github.com/gin-gonic/gin"
)

func TestHealth(t *testing.T) {
	router := newTestRouter()

	response := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodGet, "/health", nil)
	router.ServeHTTP(response, request)

	if response.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, response.Code)
	}
}

func TestTaskFlow(t *testing.T) {
	router := newTestRouter()

	createResponse := performRequest(router, http.MethodPost, "/api/tasks", `{"title":"Learn Go","description":"Build API"}`)
	if createResponse.Code != http.StatusCreated {
		t.Fatalf("expected create status %d, got %d: %s", http.StatusCreated, createResponse.Code, createResponse.Body.String())
	}

	var created model.Task
	if err := json.Unmarshal(createResponse.Body.Bytes(), &created); err != nil {
		t.Fatalf("failed to decode created task: %v", err)
	}
	if created.ID != 1 || created.Status != model.StatusTodo {
		t.Fatalf("unexpected created task: %+v", created)
	}

	statusResponse := performRequest(router, http.MethodPatch, "/api/tasks/1/status", `{"status":"in_progress"}`)
	if statusResponse.Code != http.StatusOK {
		t.Fatalf("expected update status %d, got %d: %s", http.StatusOK, statusResponse.Code, statusResponse.Body.String())
	}

	statsResponse := performRequest(router, http.MethodGet, "/api/stats", "")
	if statsResponse.Code != http.StatusOK {
		t.Fatalf("expected stats status %d, got %d: %s", http.StatusOK, statsResponse.Code, statsResponse.Body.String())
	}

	var stats model.TaskStats
	if err := json.Unmarshal(statsResponse.Body.Bytes(), &stats); err != nil {
		t.Fatalf("failed to decode stats: %v", err)
	}
	if stats.Total != 1 || stats.Todo != 0 || stats.InProgress != 1 || stats.Done != 0 {
		t.Fatalf("unexpected stats: %+v", stats)
	}

	deleteResponse := performRequest(router, http.MethodDelete, "/api/tasks/1", "")
	if deleteResponse.Code != http.StatusNoContent {
		t.Fatalf("expected delete status %d, got %d", http.StatusNoContent, deleteResponse.Code)
	}
}

func TestInvalidStatus(t *testing.T) {
	router := newTestRouter()

	performRequest(router, http.MethodPost, "/api/tasks", `{"title":"Learn Go"}`)
	response := performRequest(router, http.MethodPatch, "/api/tasks/1/status", `{"status":"blocked"}`)

	if response.Code != http.StatusBadRequest {
		t.Fatalf("expected status %d, got %d", http.StatusBadRequest, response.Code)
	}
}

func newTestRouter() http.Handler {
	gin.SetMode(gin.TestMode)
	taskRepository := repository.NewMemoryTaskRepository()
	taskService := service.NewTaskService(taskRepository, nil)
	return handler.NewRouter(taskService)
}

func performRequest(router http.Handler, method, path, body string) *httptest.ResponseRecorder {
	response := httptest.NewRecorder()
	request := httptest.NewRequest(method, path, bytes.NewBufferString(body))
	request.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(response, request)
	return response
}
