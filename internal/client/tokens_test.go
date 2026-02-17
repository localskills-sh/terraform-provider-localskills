package client

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestListUserTokens(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("expected GET, got %s", r.Method)
		}
		if r.URL.Path != "/api/user/tokens" {
			t.Errorf("expected /api/user/tokens, got %s", r.URL.Path)
		}

		resp := ApiResponse[[]ApiToken]{
			Success: true,
			Data: []ApiToken{
				{ID: "tok-1", Name: "my-token", CreatedAt: "2024-01-01T00:00:00Z"},
				{ID: "tok-2", Name: "other-token", CreatedAt: "2024-01-02T00:00:00Z"},
			},
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(resp)
	}))
	defer server.Close()

	c := NewClient(server.URL, "lsk_test123")
	tokens, err := c.ListUserTokens(context.Background())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(tokens) != 2 {
		t.Fatalf("expected 2 tokens, got %d", len(tokens))
	}
	if tokens[0].ID != "tok-1" {
		t.Errorf("expected token ID 'tok-1', got '%s'", tokens[0].ID)
	}
	if tokens[0].Name != "my-token" {
		t.Errorf("expected token name 'my-token', got '%s'", tokens[0].Name)
	}
}

func TestCreateUserToken(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("expected POST, got %s", r.Method)
		}
		if r.URL.Path != "/api/user/tokens" {
			t.Errorf("expected /api/user/tokens, got %s", r.URL.Path)
		}

		var reqBody CreateTokenRequest
		if err := json.NewDecoder(r.Body).Decode(&reqBody); err != nil {
			t.Fatalf("failed to decode request body: %v", err)
		}
		if reqBody.Name != "test-token" {
			t.Errorf("expected name 'test-token', got '%s'", reqBody.Name)
		}

		resp := ApiResponse[ApiTokenWithSecret]{
			Success: true,
			Data: ApiTokenWithSecret{
				ID:    "tok-new",
				Name:  reqBody.Name,
				Token: "lsk_secret_value",
			},
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(resp)
	}))
	defer server.Close()

	c := NewClient(server.URL, "lsk_test123")
	token, err := c.CreateUserToken(context.Background(), CreateTokenRequest{Name: "test-token"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if token.ID != "tok-new" {
		t.Errorf("expected token ID 'tok-new', got '%s'", token.ID)
	}
	if token.Token != "lsk_secret_value" {
		t.Errorf("expected token secret 'lsk_secret_value', got '%s'", token.Token)
	}
}

func TestDeleteUserToken(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodDelete {
			t.Errorf("expected DELETE, got %s", r.Method)
		}
		if r.URL.Path != "/api/user/tokens/tok-1" {
			t.Errorf("expected /api/user/tokens/tok-1, got %s", r.URL.Path)
		}

		resp := ApiResponse[struct{}]{Success: true}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(resp)
	}))
	defer server.Close()

	c := NewClient(server.URL, "lsk_test123")
	err := c.DeleteUserToken(context.Background(), "tok-1")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestListTeamTokens(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("expected GET, got %s", r.Method)
		}
		if r.URL.Path != "/api/tenants/tenant-1/tokens" {
			t.Errorf("expected /api/tenants/tenant-1/tokens, got %s", r.URL.Path)
		}

		resp := ApiResponse[[]TeamApiToken]{
			Success: true,
			Data: []TeamApiToken{
				{ID: "ttok-1", Name: "team-token", CreatedAt: "2024-01-01T00:00:00Z", CreatedByEmail: "user@test.com"},
			},
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(resp)
	}))
	defer server.Close()

	c := NewClient(server.URL, "lsk_test123")
	tokens, err := c.ListTeamTokens(context.Background(), "tenant-1")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(tokens) != 1 {
		t.Fatalf("expected 1 token, got %d", len(tokens))
	}
	if tokens[0].ID != "ttok-1" {
		t.Errorf("expected token ID 'ttok-1', got '%s'", tokens[0].ID)
	}
}

func TestCreateTeamToken(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("expected POST, got %s", r.Method)
		}
		if r.URL.Path != "/api/tenants/tenant-1/tokens" {
			t.Errorf("expected /api/tenants/tenant-1/tokens, got %s", r.URL.Path)
		}

		var reqBody CreateTeamTokenRequest
		if err := json.NewDecoder(r.Body).Decode(&reqBody); err != nil {
			t.Fatalf("failed to decode request body: %v", err)
		}
		if reqBody.Name != "team-token" {
			t.Errorf("expected name 'team-token', got '%s'", reqBody.Name)
		}

		expiresAt := "2024-04-01T00:00:00Z"
		resp := ApiResponse[TeamApiTokenWithSecret]{
			Success: true,
			Data: TeamApiTokenWithSecret{
				ID:        "ttok-new",
				Name:      reqBody.Name,
				Token:     "lsk_team_secret",
				CreatedAt: "2024-01-01T00:00:00Z",
				ExpiresAt: &expiresAt,
			},
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(resp)
	}))
	defer server.Close()

	days := 90
	c := NewClient(server.URL, "lsk_test123")
	token, err := c.CreateTeamToken(context.Background(), "tenant-1", CreateTeamTokenRequest{
		Name:          "team-token",
		ExpiresInDays: &days,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if token.ID != "ttok-new" {
		t.Errorf("expected token ID 'ttok-new', got '%s'", token.ID)
	}
	if token.Token != "lsk_team_secret" {
		t.Errorf("expected token secret 'lsk_team_secret', got '%s'", token.Token)
	}
}

func TestDeleteTeamToken(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodDelete {
			t.Errorf("expected DELETE, got %s", r.Method)
		}
		if r.URL.Path != "/api/tenants/tenant-1/tokens/ttok-1" {
			t.Errorf("expected /api/tenants/tenant-1/tokens/ttok-1, got %s", r.URL.Path)
		}

		resp := ApiResponse[struct{}]{Success: true}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(resp)
	}))
	defer server.Close()

	c := NewClient(server.URL, "lsk_test123")
	err := c.DeleteTeamToken(context.Background(), "tenant-1", "ttok-1")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestDeleteUserToken_NotFound(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)
		resp := ApiResponse[json.RawMessage]{
			Success: false,
			Error:   "token not found",
		}
		json.NewEncoder(w).Encode(resp)
	}))
	defer server.Close()

	c := NewClient(server.URL, "lsk_test123")
	err := c.DeleteUserToken(context.Background(), "nonexistent")
	if err == nil {
		t.Fatal("expected error, got nil")
	}

	if !IsNotFound(err) {
		t.Error("expected IsNotFound to return true")
	}
}
