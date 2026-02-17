package client

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestCreateInvitation(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("expected POST, got %s", r.Method)
		}
		if r.URL.Path != "/api/tenants/tenant-1/invitations" {
			t.Errorf("expected /api/tenants/tenant-1/invitations, got %s", r.URL.Path)
		}

		var reqBody CreateInvitationRequest
		if err := json.NewDecoder(r.Body).Decode(&reqBody); err != nil {
			t.Fatalf("failed to decode request body: %v", err)
		}
		if reqBody.Email != "user@example.com" {
			t.Errorf("expected email 'user@example.com', got '%s'", reqBody.Email)
		}
		if reqBody.Role != "member" {
			t.Errorf("expected role 'member', got '%s'", reqBody.Role)
		}

		w.Header().Set("Content-Type", "application/json")
		resp := ApiResponse[TenantInvitation]{
			Success: true,
			Data: TenantInvitation{
				ID:        "inv-1",
				TenantID:  "tenant-1",
				Email:     reqBody.Email,
				Role:      reqBody.Role,
				Token:     "tok_abc123",
				InvitedBy: "user-1",
				ExpiresAt: "2024-02-01T00:00:00Z",
				CreatedAt: "2024-01-01T00:00:00Z",
			},
		}
		json.NewEncoder(w).Encode(resp)
	}))
	defer server.Close()

	c := NewClient(server.URL, "lsk_test123")
	ctx := context.Background()

	invitation, err := c.CreateInvitation(ctx, "tenant-1", CreateInvitationRequest{
		Email: "user@example.com",
		Role:  "member",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if invitation.ID != "inv-1" {
		t.Errorf("expected invitation ID 'inv-1', got '%s'", invitation.ID)
	}
	if invitation.Email != "user@example.com" {
		t.Errorf("expected email 'user@example.com', got '%s'", invitation.Email)
	}
	if invitation.Token != "tok_abc123" {
		t.Errorf("expected token 'tok_abc123', got '%s'", invitation.Token)
	}
	if invitation.TenantID != "tenant-1" {
		t.Errorf("expected tenant ID 'tenant-1', got '%s'", invitation.TenantID)
	}
}

func TestListInvitations(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("expected GET, got %s", r.Method)
		}
		if r.URL.Path != "/api/tenants/tenant-1/invitations" {
			t.Errorf("expected /api/tenants/tenant-1/invitations, got %s", r.URL.Path)
		}

		w.Header().Set("Content-Type", "application/json")
		resp := ApiResponse[[]TenantInvitation]{
			Success: true,
			Data: []TenantInvitation{
				{
					ID:        "inv-1",
					TenantID:  "tenant-1",
					Email:     "user1@example.com",
					Role:      "member",
					Token:     "tok_abc123",
					InvitedBy: "user-1",
					ExpiresAt: "2024-02-01T00:00:00Z",
					CreatedAt: "2024-01-01T00:00:00Z",
				},
				{
					ID:        "inv-2",
					TenantID:  "tenant-1",
					Email:     "user2@example.com",
					Role:      "admin",
					Token:     "tok_def456",
					InvitedBy: "user-1",
					ExpiresAt: "2024-02-01T00:00:00Z",
					CreatedAt: "2024-01-01T00:00:00Z",
				},
			},
		}
		json.NewEncoder(w).Encode(resp)
	}))
	defer server.Close()

	c := NewClient(server.URL, "lsk_test123")
	ctx := context.Background()

	invitations, err := c.ListInvitations(ctx, "tenant-1")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(invitations) != 2 {
		t.Fatalf("expected 2 invitations, got %d", len(invitations))
	}
	if invitations[0].ID != "inv-1" {
		t.Errorf("expected first invitation ID 'inv-1', got '%s'", invitations[0].ID)
	}
	if invitations[0].Email != "user1@example.com" {
		t.Errorf("expected first invitation email 'user1@example.com', got '%s'", invitations[0].Email)
	}
	if invitations[1].ID != "inv-2" {
		t.Errorf("expected second invitation ID 'inv-2', got '%s'", invitations[1].ID)
	}
	if invitations[1].Role != "admin" {
		t.Errorf("expected second invitation role 'admin', got '%s'", invitations[1].Role)
	}
}

func TestCreateInvitation_Error(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		resp := ApiResponse[json.RawMessage]{
			Success: false,
			Error:   "email is required",
		}
		json.NewEncoder(w).Encode(resp)
	}))
	defer server.Close()

	c := NewClient(server.URL, "lsk_test123")
	ctx := context.Background()

	_, err := c.CreateInvitation(ctx, "tenant-1", CreateInvitationRequest{})
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

func TestListInvitations_NotFound(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)
		resp := ApiResponse[json.RawMessage]{
			Success: false,
			Error:   "tenant not found",
		}
		json.NewEncoder(w).Encode(resp)
	}))
	defer server.Close()

	c := NewClient(server.URL, "lsk_test123")
	ctx := context.Background()

	_, err := c.ListInvitations(ctx, "nonexistent")
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if !IsNotFound(err) {
		t.Error("expected IsNotFound to return true")
	}
}
