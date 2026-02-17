package client

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestGetSSOConnection(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("expected GET, got %s", r.Method)
		}
		if r.URL.Path != "/api/tenants/tenant-1/sso" {
			t.Errorf("expected /api/tenants/tenant-1/sso, got %s", r.URL.Path)
		}

		resp := ApiResponse[SsoConnection]{
			Success: true,
			Data: SsoConnection{
				ID:           "sso-1",
				TenantID:     "tenant-1",
				DisplayName:  "Corporate SSO",
				IdpEntityID:  "https://idp.example.com/entity",
				IdpSsoURL:    "https://idp.example.com/sso",
				IdpSloURL:    "https://idp.example.com/slo",
				SpEntityID:   "https://localskills.sh/sp",
				SpAcsURL:     "https://localskills.sh/acs",
				DefaultRole:  "member",
				EmailDomains: []string{"example.com"},
				Enabled:      true,
				RequireSso:   false,
				MetadataURL:  "https://idp.example.com/metadata",
				CreatedAt:    "2024-01-01T00:00:00Z",
				UpdatedAt:    "2024-01-01T00:00:00Z",
			},
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(resp)
	}))
	defer server.Close()

	c := NewClient(server.URL, "lsk_test123")
	conn, err := c.GetSSOConnection(context.Background(), "tenant-1")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if conn.ID != "sso-1" {
		t.Errorf("expected SSO ID 'sso-1', got '%s'", conn.ID)
	}
	if conn.DisplayName != "Corporate SSO" {
		t.Errorf("expected display name 'Corporate SSO', got '%s'", conn.DisplayName)
	}
	if !conn.Enabled {
		t.Error("expected SSO to be enabled")
	}
	if len(conn.EmailDomains) != 1 || conn.EmailDomains[0] != "example.com" {
		t.Errorf("unexpected email domains: %v", conn.EmailDomains)
	}
}

func TestGetSSOConnection_NotFound(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)
		resp := ApiResponse[json.RawMessage]{
			Success: false,
			Error:   "SSO connection not found",
		}
		json.NewEncoder(w).Encode(resp)
	}))
	defer server.Close()

	c := NewClient(server.URL, "lsk_test123")
	_, err := c.GetSSOConnection(context.Background(), "nonexistent")
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if !IsNotFound(err) {
		t.Error("expected IsNotFound to return true")
	}
}

func TestUpdateSSOConnection(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPatch {
			t.Errorf("expected PATCH, got %s", r.Method)
		}
		if r.URL.Path != "/api/tenants/tenant-1/sso" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}

		var reqBody UpdateSsoRequest
		if err := json.NewDecoder(r.Body).Decode(&reqBody); err != nil {
			t.Fatalf("failed to decode request body: %v", err)
		}
		if reqBody.DisplayName == nil || *reqBody.DisplayName != "Updated SSO" {
			t.Errorf("expected display name 'Updated SSO', got %v", reqBody.DisplayName)
		}
		if reqBody.Enabled == nil || !*reqBody.Enabled {
			t.Errorf("expected enabled=true")
		}

		resp := ApiResponse[SsoConnection]{
			Success: true,
			Data: SsoConnection{
				ID:           "sso-1",
				TenantID:     "tenant-1",
				DisplayName:  *reqBody.DisplayName,
				IdpEntityID:  "https://idp.example.com/entity",
				IdpSsoURL:    "https://idp.example.com/sso",
				IdpSloURL:    "https://idp.example.com/slo",
				SpEntityID:   "https://localskills.sh/sp",
				SpAcsURL:     "https://localskills.sh/acs",
				DefaultRole:  "member",
				EmailDomains: []string{"example.com"},
				Enabled:      *reqBody.Enabled,
				RequireSso:   false,
				MetadataURL:  "https://idp.example.com/metadata",
				CreatedAt:    "2024-01-01T00:00:00Z",
				UpdatedAt:    "2024-01-02T00:00:00Z",
			},
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(resp)
	}))
	defer server.Close()

	c := NewClient(server.URL, "lsk_test123")
	displayName := "Updated SSO"
	enabled := true
	conn, err := c.UpdateSSOConnection(context.Background(), "tenant-1", UpdateSsoRequest{
		DisplayName: &displayName,
		Enabled:     &enabled,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if conn.DisplayName != "Updated SSO" {
		t.Errorf("expected display name 'Updated SSO', got '%s'", conn.DisplayName)
	}
	if !conn.Enabled {
		t.Error("expected SSO to be enabled")
	}
	if conn.UpdatedAt != "2024-01-02T00:00:00Z" {
		t.Errorf("expected updated timestamp, got '%s'", conn.UpdatedAt)
	}
}
