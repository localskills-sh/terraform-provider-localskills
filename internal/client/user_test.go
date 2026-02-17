package client

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestGetUserProfile(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("expected GET, got %s", r.Method)
		}
		if r.URL.Path != "/api/user/profile" {
			t.Errorf("expected /api/user/profile, got %s", r.URL.Path)
		}

		username := "testuser"
		name := "Test User"
		image := "https://example.com/avatar.png"
		bio := "A test user"
		resp := ApiResponse[UserProfile]{
			Success: true,
			Data: UserProfile{
				ID:       "user-1",
				Username: &username,
				Name:     &name,
				Email:    "test@example.com",
				Image:    &image,
				Bio:      &bio,
			},
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(resp)
	}))
	defer server.Close()

	c := NewClient(server.URL, "lsk_test123")
	profile, err := c.GetUserProfile(context.Background())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if profile.ID != "user-1" {
		t.Errorf("expected ID 'user-1', got '%s'", profile.ID)
	}
	if profile.Email != "test@example.com" {
		t.Errorf("expected email 'test@example.com', got '%s'", profile.Email)
	}
	if profile.Username == nil || *profile.Username != "testuser" {
		t.Errorf("expected username 'testuser', got %v", profile.Username)
	}
}

func TestListUserAuditLog(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("expected GET, got %s", r.Method)
		}
		if r.URL.Path != "/api/user/audit-log" {
			t.Errorf("expected /api/user/audit-log, got %s", r.URL.Path)
		}

		actorID := "user-1"
		actorName := "Test User"
		resp := ApiResponse[AuditLogResponse]{
			Success: true,
			Data: AuditLogResponse{
				Entries: []AuditLogEntry{
					{
						ID:           "log-1",
						Action:       "skill.create",
						ActorID:      &actorID,
						ActorName:    &actorName,
						ResourceType: "skill",
						ResourceID:   "skill-1",
						Metadata:     "{}",
						CreatedAt:    "2024-01-01T00:00:00Z",
					},
				},
				Total: 1,
				Page:  1,
			},
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(resp)
	}))
	defer server.Close()

	c := NewClient(server.URL, "lsk_test123")
	result, err := c.ListUserAuditLog(context.Background(), map[string]string{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(result.Entries) != 1 {
		t.Fatalf("expected 1 entry, got %d", len(result.Entries))
	}
	if result.Entries[0].ID != "log-1" {
		t.Errorf("expected entry ID 'log-1', got '%s'", result.Entries[0].ID)
	}
	if result.Entries[0].Action != "skill.create" {
		t.Errorf("expected action 'skill.create', got '%s'", result.Entries[0].Action)
	}
}

func TestListUserAuditLog_WithParams(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Query().Get("page") != "2" {
			t.Errorf("expected page=2, got %s", r.URL.Query().Get("page"))
		}
		if r.URL.Query().Get("action") != "skill.create" {
			t.Errorf("expected action=skill.create, got %s", r.URL.Query().Get("action"))
		}

		resp := ApiResponse[AuditLogResponse]{
			Success: true,
			Data: AuditLogResponse{
				Entries: []AuditLogEntry{},
				Total:   0,
				Page:    2,
			},
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(resp)
	}))
	defer server.Close()

	c := NewClient(server.URL, "lsk_test123")
	result, err := c.ListUserAuditLog(context.Background(), map[string]string{
		"page":   "2",
		"action": "skill.create",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if result.Page != 2 {
		t.Errorf("expected page 2, got %d", result.Page)
	}
}

func TestListTeamAuditLog(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("expected GET, got %s", r.Method)
		}
		if r.URL.Path != "/api/tenants/tenant-1/audit-log" {
			t.Errorf("expected /api/tenants/tenant-1/audit-log, got %s", r.URL.Path)
		}

		actorID := "user-1"
		actorName := "Team Member"
		resp := ApiResponse[AuditLogResponse]{
			Success: true,
			Data: AuditLogResponse{
				Entries: []AuditLogEntry{
					{
						ID:           "log-2",
						Action:       "team.member.add",
						ActorID:      &actorID,
						ActorName:    &actorName,
						ResourceType: "team",
						ResourceID:   "tenant-1",
						Metadata:     "{}",
						CreatedAt:    "2024-01-01T00:00:00Z",
					},
				},
				Total: 1,
				Page:  1,
			},
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(resp)
	}))
	defer server.Close()

	c := NewClient(server.URL, "lsk_test123")
	result, err := c.ListTeamAuditLog(context.Background(), "tenant-1", map[string]string{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(result.Entries) != 1 {
		t.Fatalf("expected 1 entry, got %d", len(result.Entries))
	}
	if result.Entries[0].ID != "log-2" {
		t.Errorf("expected entry ID 'log-2', got '%s'", result.Entries[0].ID)
	}
	if result.Entries[0].Action != "team.member.add" {
		t.Errorf("expected action 'team.member.add', got '%s'", result.Entries[0].Action)
	}
}
