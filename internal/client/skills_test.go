package client

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestCreateSkill(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("expected POST, got %s", r.Method)
		}
		if r.URL.Path != "/api/skills" {
			t.Errorf("expected /api/skills, got %s", r.URL.Path)
		}

		var req CreateSkillRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			t.Fatalf("failed to decode request: %v", err)
		}
		if req.Name != "test-skill" {
			t.Errorf("expected name 'test-skill', got '%s'", req.Name)
		}
		if req.Type != "skill" {
			t.Errorf("expected type 'skill', got '%s'", req.Type)
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(ApiResponse[Skill]{
			Success: true,
			Data: Skill{
				ID:   "skill-123",
				Name: req.Name,
				Type: req.Type,
				Slug: "test-skill",
			},
		})
	}))
	defer server.Close()

	c := NewClient(server.URL, "lsk_test123")
	skill, err := c.CreateSkill(context.Background(), CreateSkillRequest{
		Name:       "test-skill",
		Type:       "skill",
		Visibility: "private",
		Content:    "# Test",
		TenantID:   "tenant-1",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if skill.ID != "skill-123" {
		t.Errorf("expected ID 'skill-123', got '%s'", skill.ID)
	}
	if skill.Name != "test-skill" {
		t.Errorf("expected name 'test-skill', got '%s'", skill.Name)
	}
}

func TestCreateSkill_Error(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(ApiResponse[json.RawMessage]{
			Success: false,
			Error:   "name is required",
		})
	}))
	defer server.Close()

	c := NewClient(server.URL, "lsk_test123")
	_, err := c.CreateSkill(context.Background(), CreateSkillRequest{})
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

func TestGetSkill(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("expected GET, got %s", r.Method)
		}
		if r.URL.Path != "/api/skills/skill-123" {
			t.Errorf("expected /api/skills/skill-123, got %s", r.URL.Path)
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(ApiResponse[SkillWithVersion]{
			Success: true,
			Data: SkillWithVersion{
				Skill: Skill{
					ID:   "skill-123",
					Name: "test-skill",
				},
				CurrentVersionInfo: &SkillVersion{
					ID:      "ver-1",
					Version: 1,
				},
			},
		})
	}))
	defer server.Close()

	c := NewClient(server.URL, "lsk_test123")
	skill, err := c.GetSkill(context.Background(), "skill-123")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if skill.ID != "skill-123" {
		t.Errorf("expected ID 'skill-123', got '%s'", skill.ID)
	}
	if skill.CurrentVersionInfo == nil {
		t.Fatal("expected CurrentVersionInfo to be set")
	}
	if skill.CurrentVersionInfo.Version != 1 {
		t.Errorf("expected version 1, got %d", skill.CurrentVersionInfo.Version)
	}
}

func TestGetSkill_NotFound(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(ApiResponse[json.RawMessage]{
			Success: false,
			Error:   "skill not found",
		})
	}))
	defer server.Close()

	c := NewClient(server.URL, "lsk_test123")
	_, err := c.GetSkill(context.Background(), "nonexistent")
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if !IsNotFound(err) {
		t.Error("expected IsNotFound to return true")
	}
}

func TestUpdateSkill(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPut {
			t.Errorf("expected PUT, got %s", r.Method)
		}
		if r.URL.Path != "/api/skills/skill-123" {
			t.Errorf("expected /api/skills/skill-123, got %s", r.URL.Path)
		}

		var req UpdateSkillRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			t.Fatalf("failed to decode request: %v", err)
		}
		if req.Name == nil || *req.Name != "updated-skill" {
			t.Error("expected name to be 'updated-skill'")
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(ApiResponse[Skill]{
			Success: true,
			Data: Skill{
				ID:   "skill-123",
				Name: "updated-skill",
			},
		})
	}))
	defer server.Close()

	c := NewClient(server.URL, "lsk_test123")
	name := "updated-skill"
	skill, err := c.UpdateSkill(context.Background(), "skill-123", UpdateSkillRequest{
		Name: &name,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if skill.Name != "updated-skill" {
		t.Errorf("expected name 'updated-skill', got '%s'", skill.Name)
	}
}

func TestDeleteSkill(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodDelete {
			t.Errorf("expected DELETE, got %s", r.Method)
		}
		if r.URL.Path != "/api/skills/skill-123" {
			t.Errorf("expected /api/skills/skill-123, got %s", r.URL.Path)
		}
		w.WriteHeader(http.StatusNoContent)
	}))
	defer server.Close()

	c := NewClient(server.URL, "lsk_test123")
	err := c.DeleteSkill(context.Background(), "skill-123")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestDeleteSkill_NotFound(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	}))
	defer server.Close()

	c := NewClient(server.URL, "lsk_test123")
	err := c.DeleteSkill(context.Background(), "nonexistent")
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

func TestListSkills(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("expected GET, got %s", r.Method)
		}
		if !hasPrefix(r.URL.Path, "/api/skills") {
			t.Errorf("expected path starting with /api/skills, got %s", r.URL.Path)
		}
		if r.URL.Query().Get("tenant_id") != "tenant-1" {
			t.Errorf("expected tenant_id=tenant-1, got %s", r.URL.Query().Get("tenant_id"))
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(ApiResponse[[]Skill]{
			Success: true,
			Data: []Skill{
				{ID: "skill-1", Name: "skill-one"},
				{ID: "skill-2", Name: "skill-two"},
			},
		})
	}))
	defer server.Close()

	c := NewClient(server.URL, "lsk_test123")
	skills, err := c.ListSkills(context.Background(), map[string]string{
		"tenant_id": "tenant-1",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(skills) != 2 {
		t.Fatalf("expected 2 skills, got %d", len(skills))
	}
	if skills[0].ID != "skill-1" {
		t.Errorf("expected first skill ID 'skill-1', got '%s'", skills[0].ID)
	}
}

func TestListSkills_NoParams(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.RawQuery != "" {
			t.Errorf("expected no query params, got %s", r.URL.RawQuery)
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(ApiResponse[[]Skill]{
			Success: true,
			Data:    []Skill{},
		})
	}))
	defer server.Close()

	c := NewClient(server.URL, "lsk_test123")
	skills, err := c.ListSkills(context.Background(), nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(skills) != 0 {
		t.Errorf("expected 0 skills, got %d", len(skills))
	}
}

func hasPrefix(s, prefix string) bool {
	return len(s) >= len(prefix) && s[:len(prefix)] == prefix
}
