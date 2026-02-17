package client

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestListOIDCPolicies(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("expected GET, got %s", r.Method)
		}
		if r.URL.Path != "/api/tenants/tenant-1/oidc-policies" {
			t.Errorf("expected /api/tenants/tenant-1/oidc-policies, got %s", r.URL.Path)
		}

		resp := ApiResponse[[]OidcTrustPolicy]{
			Success: true,
			Data: []OidcTrustPolicy{
				{
					ID:         "policy-1",
					TenantID:   "tenant-1",
					Name:       "github-deploy",
					Provider:   "github",
					Repository: "org/repo",
					RefFilter:  "*",
					Enabled:    true,
					CreatedBy:  "user-1",
					CreatedAt:  "2024-01-01T00:00:00Z",
					UpdatedAt:  "2024-01-01T00:00:00Z",
				},
			},
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(resp)
	}))
	defer server.Close()

	c := NewClient(server.URL, "lsk_test123")
	policies, err := c.ListOIDCPolicies(context.Background(), "tenant-1")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(policies) != 1 {
		t.Fatalf("expected 1 policy, got %d", len(policies))
	}
	if policies[0].ID != "policy-1" {
		t.Errorf("expected policy ID 'policy-1', got '%s'", policies[0].ID)
	}
	if policies[0].Provider != "github" {
		t.Errorf("expected provider 'github', got '%s'", policies[0].Provider)
	}
}

func TestCreateOIDCPolicy(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("expected POST, got %s", r.Method)
		}
		if r.URL.Path != "/api/tenants/tenant-1/oidc-policies" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}

		var reqBody CreateOidcPolicyRequest
		if err := json.NewDecoder(r.Body).Decode(&reqBody); err != nil {
			t.Fatalf("failed to decode request body: %v", err)
		}
		if reqBody.Name != "github-deploy" {
			t.Errorf("expected name 'github-deploy', got '%s'", reqBody.Name)
		}
		if reqBody.Provider != "github" {
			t.Errorf("expected provider 'github', got '%s'", reqBody.Provider)
		}

		resp := ApiResponse[OidcTrustPolicy]{
			Success: true,
			Data: OidcTrustPolicy{
				ID:         "policy-new",
				TenantID:   "tenant-1",
				Name:       reqBody.Name,
				Provider:   reqBody.Provider,
				Repository: reqBody.Repository,
				RefFilter:  "*",
				Enabled:    true,
				CreatedBy:  "user-1",
				CreatedAt:  "2024-01-01T00:00:00Z",
				UpdatedAt:  "2024-01-01T00:00:00Z",
			},
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(resp)
	}))
	defer server.Close()

	c := NewClient(server.URL, "lsk_test123")
	policy, err := c.CreateOIDCPolicy(context.Background(), "tenant-1", CreateOidcPolicyRequest{
		Name:       "github-deploy",
		Provider:   "github",
		Repository: "org/repo",
		Enabled:    true,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if policy.ID != "policy-new" {
		t.Errorf("expected policy ID 'policy-new', got '%s'", policy.ID)
	}
	if policy.Name != "github-deploy" {
		t.Errorf("expected name 'github-deploy', got '%s'", policy.Name)
	}
}

func TestUpdateOIDCPolicy(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPatch {
			t.Errorf("expected PATCH, got %s", r.Method)
		}
		if r.URL.Path != "/api/tenants/tenant-1/oidc-policies/policy-1" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}

		var reqBody UpdateOidcPolicyRequest
		if err := json.NewDecoder(r.Body).Decode(&reqBody); err != nil {
			t.Fatalf("failed to decode request body: %v", err)
		}
		if reqBody.Name == nil || *reqBody.Name != "updated-name" {
			t.Errorf("expected name 'updated-name', got %v", reqBody.Name)
		}

		resp := ApiResponse[OidcTrustPolicy]{
			Success: true,
			Data: OidcTrustPolicy{
				ID:         "policy-1",
				TenantID:   "tenant-1",
				Name:       *reqBody.Name,
				Provider:   "github",
				Repository: "org/repo",
				RefFilter:  "*",
				Enabled:    true,
				CreatedBy:  "user-1",
				CreatedAt:  "2024-01-01T00:00:00Z",
				UpdatedAt:  "2024-01-02T00:00:00Z",
			},
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(resp)
	}))
	defer server.Close()

	c := NewClient(server.URL, "lsk_test123")
	name := "updated-name"
	policy, err := c.UpdateOIDCPolicy(context.Background(), "tenant-1", "policy-1", UpdateOidcPolicyRequest{
		Name: &name,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if policy.Name != "updated-name" {
		t.Errorf("expected name 'updated-name', got '%s'", policy.Name)
	}
	if policy.UpdatedAt != "2024-01-02T00:00:00Z" {
		t.Errorf("expected updated timestamp, got '%s'", policy.UpdatedAt)
	}
}

func TestDeleteOIDCPolicy(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodDelete {
			t.Errorf("expected DELETE, got %s", r.Method)
		}
		if r.URL.Path != "/api/tenants/tenant-1/oidc-policies/policy-1" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}

		resp := ApiResponse[struct{}]{
			Success: true,
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(resp)
	}))
	defer server.Close()

	c := NewClient(server.URL, "lsk_test123")
	err := c.DeleteOIDCPolicy(context.Background(), "tenant-1", "policy-1")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestDeleteOIDCPolicy_NotFound(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)
		resp := ApiResponse[json.RawMessage]{
			Success: false,
			Error:   "policy not found",
		}
		json.NewEncoder(w).Encode(resp)
	}))
	defer server.Close()

	c := NewClient(server.URL, "lsk_test123")
	err := c.DeleteOIDCPolicy(context.Background(), "tenant-1", "nonexistent")
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if !IsNotFound(err) {
		t.Error("expected IsNotFound to return true")
	}
}
