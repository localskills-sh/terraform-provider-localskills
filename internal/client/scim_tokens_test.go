package client

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestListSCIMTokens(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("expected GET, got %s", r.Method)
		}
		if r.URL.Path != "/api/tenants/tenant-1/scim-tokens" {
			t.Errorf("expected /api/tenants/tenant-1/scim-tokens, got %s", r.URL.Path)
		}

		resp := ApiResponse[[]ScimToken]{
			Success: true,
			Data: []ScimToken{
				{ID: "scim-1", Name: "scim-token", CreatedAt: "2024-01-01T00:00:00Z"},
			},
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(resp)
	}))
	defer server.Close()

	c := NewClient(server.URL, "lsk_test123")
	tokens, err := c.ListSCIMTokens(context.Background(), "tenant-1")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(tokens) != 1 {
		t.Fatalf("expected 1 token, got %d", len(tokens))
	}
	if tokens[0].ID != "scim-1" {
		t.Errorf("expected token ID 'scim-1', got '%s'", tokens[0].ID)
	}
	if tokens[0].Name != "scim-token" {
		t.Errorf("expected token name 'scim-token', got '%s'", tokens[0].Name)
	}
}

func TestCreateSCIMToken(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("expected POST, got %s", r.Method)
		}
		if r.URL.Path != "/api/tenants/tenant-1/scim-tokens" {
			t.Errorf("expected /api/tenants/tenant-1/scim-tokens, got %s", r.URL.Path)
		}

		var reqBody CreateScimTokenRequest
		if err := json.NewDecoder(r.Body).Decode(&reqBody); err != nil {
			t.Fatalf("failed to decode request body: %v", err)
		}
		if reqBody.Name != "scim-token" {
			t.Errorf("expected name 'scim-token', got '%s'", reqBody.Name)
		}

		expiresAt := "2024-04-01T00:00:00Z"
		resp := ApiResponse[ScimTokenWithSecret]{
			Success: true,
			Data: ScimTokenWithSecret{
				ID:        "scim-new",
				Name:      reqBody.Name,
				Token:     "lsk_scim_secret",
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
	token, err := c.CreateSCIMToken(context.Background(), "tenant-1", CreateScimTokenRequest{
		Name:          "scim-token",
		ExpiresInDays: &days,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if token.ID != "scim-new" {
		t.Errorf("expected token ID 'scim-new', got '%s'", token.ID)
	}
	if token.Token != "lsk_scim_secret" {
		t.Errorf("expected token secret 'lsk_scim_secret', got '%s'", token.Token)
	}
	if token.ExpiresAt == nil || *token.ExpiresAt != "2024-04-01T00:00:00Z" {
		t.Errorf("unexpected expires_at value")
	}
}

func TestDeleteSCIMToken(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodDelete {
			t.Errorf("expected DELETE, got %s", r.Method)
		}
		if r.URL.Path != "/api/tenants/tenant-1/scim-tokens/scim-1" {
			t.Errorf("expected /api/tenants/tenant-1/scim-tokens/scim-1, got %s", r.URL.Path)
		}

		resp := ApiResponse[struct{}]{Success: true}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(resp)
	}))
	defer server.Close()

	c := NewClient(server.URL, "lsk_test123")
	err := c.DeleteSCIMToken(context.Background(), "tenant-1", "scim-1")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestDeleteSCIMToken_NotFound(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)
		resp := ApiResponse[json.RawMessage]{
			Success: false,
			Error:   "SCIM token not found",
		}
		json.NewEncoder(w).Encode(resp)
	}))
	defer server.Close()

	c := NewClient(server.URL, "lsk_test123")
	err := c.DeleteSCIMToken(context.Background(), "tenant-1", "nonexistent")
	if err == nil {
		t.Fatal("expected error, got nil")
	}

	if !IsNotFound(err) {
		t.Error("expected IsNotFound to return true")
	}
}
