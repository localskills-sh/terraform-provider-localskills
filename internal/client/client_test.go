package client

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"sync/atomic"
	"testing"
	"time"
)

func TestDoJSON_HappyPath(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("expected GET, got %s", r.Method)
		}
		if r.URL.Path != "/api/v1/skills" {
			t.Errorf("expected /api/v1/skills, got %s", r.URL.Path)
		}
		if r.Header.Get("Authorization") != "Bearer lsk_test123" {
			t.Errorf("unexpected auth header: %s", r.Header.Get("Authorization"))
		}
		if r.Header.Get("Content-Type") != "application/json" {
			t.Errorf("unexpected content type: %s", r.Header.Get("Content-Type"))
		}

		resp := ApiResponse[[]Skill]{
			Success: true,
			Data: []Skill{
				{ID: "skill-1", Name: "test-skill", Slug: "test-skill"},
			},
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(resp)
	}))
	defer server.Close()

	c := NewClient(server.URL, "lsk_test123")
	ctx := context.Background()

	skills, err := DoJSON[[]Skill](c, ctx, http.MethodGet, "/api/v1/skills", nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(*skills) != 1 {
		t.Fatalf("expected 1 skill, got %d", len(*skills))
	}
	if (*skills)[0].ID != "skill-1" {
		t.Errorf("expected skill ID 'skill-1', got '%s'", (*skills)[0].ID)
	}
	if (*skills)[0].Name != "test-skill" {
		t.Errorf("expected skill name 'test-skill', got '%s'", (*skills)[0].Name)
	}
}

func TestDoJSON_ErrorResponse(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		resp := ApiResponse[json.RawMessage]{
			Success: false,
			Error:   "invalid request: name is required",
		}
		json.NewEncoder(w).Encode(resp)
	}))
	defer server.Close()

	c := NewClient(server.URL, "lsk_test123")
	ctx := context.Background()

	_, err := DoJSON[Skill](c, ctx, http.MethodPost, "/api/v1/skills", nil)
	if err == nil {
		t.Fatal("expected error, got nil")
	}

	apiErr, ok := err.(*ApiError)
	if !ok {
		t.Fatalf("expected *ApiError, got %T", err)
	}
	if apiErr.StatusCode != 400 {
		t.Errorf("expected status 400, got %d", apiErr.StatusCode)
	}
	if apiErr.Message != "invalid request: name is required" {
		t.Errorf("unexpected message: %s", apiErr.Message)
	}
}

func TestDoJSON_NotFound(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)
		resp := ApiResponse[json.RawMessage]{
			Success: false,
			Error:   "skill not found",
		}
		json.NewEncoder(w).Encode(resp)
	}))
	defer server.Close()

	c := NewClient(server.URL, "lsk_test123")
	ctx := context.Background()

	_, err := DoJSON[Skill](c, ctx, http.MethodGet, "/api/v1/skills/nonexistent", nil)
	if err == nil {
		t.Fatal("expected error, got nil")
	}

	if !IsNotFound(err) {
		t.Error("expected IsNotFound to return true")
	}
	if IsUnauthorized(err) {
		t.Error("expected IsUnauthorized to return false")
	}
}

func TestDoJSON_RetryOn429(t *testing.T) {
	var attempts atomic.Int32

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		attempt := attempts.Add(1)
		if attempt <= 2 {
			w.WriteHeader(http.StatusTooManyRequests)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		resp := ApiResponse[Skill]{
			Success: true,
			Data:    Skill{ID: "skill-1", Name: "test-skill"},
		}
		json.NewEncoder(w).Encode(resp)
	}))
	defer server.Close()

	c := NewClient(server.URL, "lsk_test123")
	// Use a short timeout context that's still long enough for retries
	ctx := context.Background()

	skill, err := DoJSON[Skill](c, ctx, http.MethodGet, "/api/v1/skills/1", nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if skill.ID != "skill-1" {
		t.Errorf("expected skill ID 'skill-1', got '%s'", skill.ID)
	}

	if attempts.Load() != 3 {
		t.Errorf("expected 3 attempts, got %d", attempts.Load())
	}
}

func TestDoJSON_NetworkTimeout(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(2 * time.Second)
	}))
	defer server.Close()

	c := NewClient(server.URL, "lsk_test123")
	c.HTTPClient.Timeout = 100 * time.Millisecond

	ctx := context.Background()
	_, err := DoJSON[Skill](c, ctx, http.MethodGet, "/api/v1/skills/1", nil)
	if err == nil {
		t.Fatal("expected timeout error, got nil")
	}
}

func TestDoJSON_SuccessFalse(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		resp := ApiResponse[json.RawMessage]{
			Success: false,
			Error:   "operation failed",
		}
		json.NewEncoder(w).Encode(resp)
	}))
	defer server.Close()

	c := NewClient(server.URL, "lsk_test123")
	ctx := context.Background()

	_, err := DoJSON[Skill](c, ctx, http.MethodGet, "/api/v1/skills/1", nil)
	if err == nil {
		t.Fatal("expected error, got nil")
	}

	apiErr, ok := err.(*ApiError)
	if !ok {
		t.Fatalf("expected *ApiError, got %T", err)
	}
	if apiErr.Message != "operation failed" {
		t.Errorf("unexpected message: %s", apiErr.Message)
	}
}

func TestDoJSON_PostWithBody(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("expected POST, got %s", r.Method)
		}

		var reqBody CreateSkillRequest
		if err := json.NewDecoder(r.Body).Decode(&reqBody); err != nil {
			t.Fatalf("failed to decode request body: %v", err)
		}
		if reqBody.Name != "my-skill" {
			t.Errorf("expected name 'my-skill', got '%s'", reqBody.Name)
		}

		w.Header().Set("Content-Type", "application/json")
		resp := ApiResponse[Skill]{
			Success: true,
			Data:    Skill{ID: "new-skill-1", Name: reqBody.Name},
		}
		json.NewEncoder(w).Encode(resp)
	}))
	defer server.Close()

	c := NewClient(server.URL, "lsk_test123")
	ctx := context.Background()

	body := CreateSkillRequest{
		Name:       "my-skill",
		Type:       "prompt",
		Visibility: "private",
		Content:    "test content",
		TenantID:   "tenant-1",
	}

	skill, err := DoJSON[Skill](c, ctx, http.MethodPost, "/api/v1/skills", body)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if skill.ID != "new-skill-1" {
		t.Errorf("expected skill ID 'new-skill-1', got '%s'", skill.ID)
	}
}
