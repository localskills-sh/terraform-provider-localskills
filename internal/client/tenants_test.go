package client

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestCreateTenant(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("expected POST, got %s", r.Method)
		}
		if r.URL.Path != "/api/tenants" {
			t.Errorf("expected /api/tenants, got %s", r.URL.Path)
		}

		var reqBody CreateTenantRequest
		if err := json.NewDecoder(r.Body).Decode(&reqBody); err != nil {
			t.Fatalf("failed to decode request body: %v", err)
		}
		if reqBody.Name != "my-team" {
			t.Errorf("expected name 'my-team', got '%s'", reqBody.Name)
		}

		w.Header().Set("Content-Type", "application/json")
		resp := ApiResponse[Tenant]{
			Success: true,
			Data: Tenant{
				ID:          "tenant-1",
				Name:        reqBody.Name,
				Slug:        "my-team",
				Description: "",
				CreatedAt:   "2024-01-01T00:00:00Z",
				UpdatedAt:   "2024-01-01T00:00:00Z",
			},
		}
		json.NewEncoder(w).Encode(resp)
	}))
	defer server.Close()

	c := NewClient(server.URL, "lsk_test123")
	ctx := context.Background()

	tenant, err := c.CreateTenant(ctx, CreateTenantRequest{Name: "my-team"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if tenant.ID != "tenant-1" {
		t.Errorf("expected tenant ID 'tenant-1', got '%s'", tenant.ID)
	}
	if tenant.Name != "my-team" {
		t.Errorf("expected tenant name 'my-team', got '%s'", tenant.Name)
	}
	if tenant.Slug != "my-team" {
		t.Errorf("expected tenant slug 'my-team', got '%s'", tenant.Slug)
	}
}

func TestListTenants(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("expected GET, got %s", r.Method)
		}
		if r.URL.Path != "/api/tenants" {
			t.Errorf("expected /api/tenants, got %s", r.URL.Path)
		}

		w.Header().Set("Content-Type", "application/json")
		resp := ApiResponse[[]TenantWithRole]{
			Success: true,
			Data: []TenantWithRole{
				{
					ID:          "tenant-1",
					Name:        "Team One",
					Slug:        "team-one",
					Description: "First team",
					Role:        "owner",
					CreatedAt:   "2024-01-01T00:00:00Z",
					UpdatedAt:   "2024-01-01T00:00:00Z",
				},
				{
					ID:          "tenant-2",
					Name:        "Team Two",
					Slug:        "team-two",
					Description: "Second team",
					Role:        "member",
					CreatedAt:   "2024-01-02T00:00:00Z",
					UpdatedAt:   "2024-01-02T00:00:00Z",
				},
			},
		}
		json.NewEncoder(w).Encode(resp)
	}))
	defer server.Close()

	c := NewClient(server.URL, "lsk_test123")
	ctx := context.Background()

	tenants, err := c.ListTenants(ctx)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(tenants) != 2 {
		t.Fatalf("expected 2 tenants, got %d", len(tenants))
	}
	if tenants[0].ID != "tenant-1" {
		t.Errorf("expected first tenant ID 'tenant-1', got '%s'", tenants[0].ID)
	}
	if tenants[0].Role != "owner" {
		t.Errorf("expected first tenant role 'owner', got '%s'", tenants[0].Role)
	}
	if tenants[1].ID != "tenant-2" {
		t.Errorf("expected second tenant ID 'tenant-2', got '%s'", tenants[1].ID)
	}
}

func TestUpdateTenant(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPatch {
			t.Errorf("expected PATCH, got %s", r.Method)
		}
		if r.URL.Path != "/api/tenants/tenant-1" {
			t.Errorf("expected /api/tenants/tenant-1, got %s", r.URL.Path)
		}

		var reqBody UpdateTenantRequest
		if err := json.NewDecoder(r.Body).Decode(&reqBody); err != nil {
			t.Fatalf("failed to decode request body: %v", err)
		}
		if reqBody.Name == nil || *reqBody.Name != "Updated Team" {
			t.Errorf("expected name 'Updated Team', got %v", reqBody.Name)
		}

		w.Header().Set("Content-Type", "application/json")
		resp := ApiResponse[Tenant]{
			Success: true,
			Data: Tenant{
				ID:          "tenant-1",
				Name:        *reqBody.Name,
				Slug:        "updated-team",
				Description: "updated description",
				CreatedAt:   "2024-01-01T00:00:00Z",
				UpdatedAt:   "2024-01-02T00:00:00Z",
			},
		}
		json.NewEncoder(w).Encode(resp)
	}))
	defer server.Close()

	c := NewClient(server.URL, "lsk_test123")
	ctx := context.Background()

	name := "Updated Team"
	tenant, err := c.UpdateTenant(ctx, "tenant-1", UpdateTenantRequest{Name: &name})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if tenant.ID != "tenant-1" {
		t.Errorf("expected tenant ID 'tenant-1', got '%s'", tenant.ID)
	}
	if tenant.Name != "Updated Team" {
		t.Errorf("expected tenant name 'Updated Team', got '%s'", tenant.Name)
	}
	if tenant.Slug != "updated-team" {
		t.Errorf("expected tenant slug 'updated-team', got '%s'", tenant.Slug)
	}
}

func TestCreateTenant_Error(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		resp := ApiResponse[json.RawMessage]{
			Success: false,
			Error:   "name is required",
		}
		json.NewEncoder(w).Encode(resp)
	}))
	defer server.Close()

	c := NewClient(server.URL, "lsk_test123")
	ctx := context.Background()

	_, err := c.CreateTenant(ctx, CreateTenantRequest{})
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
}
