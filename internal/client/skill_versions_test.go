package client

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestCreateSkillVersion(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("expected POST, got %s", r.Method)
		}
		if r.URL.Path != "/api/skills/skill-123/versions" {
			t.Errorf("expected /api/skills/skill-123/versions, got %s", r.URL.Path)
		}

		var req CreateSkillVersionRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			t.Fatalf("failed to decode request: %v", err)
		}
		if req.Content != "# Updated content" {
			t.Errorf("expected content '# Updated content', got '%s'", req.Content)
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(ApiResponse[SkillVersion]{
			Success: true,
			Data: SkillVersion{
				ID:      "ver-2",
				SkillID: "skill-123",
				Version: 2,
				Semver:  "1.1.0",
			},
		})
	}))
	defer server.Close()

	c := NewClient(server.URL, "lsk_test123")
	ver, err := c.CreateSkillVersion(context.Background(), "skill-123", CreateSkillVersionRequest{
		Content: "# Updated content",
		Message: "update content",
		Bump:    "minor",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if ver.ID != "ver-2" {
		t.Errorf("expected ID 'ver-2', got '%s'", ver.ID)
	}
	if ver.Version != 2 {
		t.Errorf("expected version 2, got %d", ver.Version)
	}
}

func TestCreateSkillVersion_Error(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(ApiResponse[json.RawMessage]{
			Success: false,
			Error:   "content is required",
		})
	}))
	defer server.Close()

	c := NewClient(server.URL, "lsk_test123")
	_, err := c.CreateSkillVersion(context.Background(), "skill-123", CreateSkillVersionRequest{})
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

func TestListSkillVersions(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("expected GET, got %s", r.Method)
		}
		if r.URL.Path != "/api/skills/skill-123/versions" {
			t.Errorf("expected /api/skills/skill-123/versions, got %s", r.URL.Path)
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(ApiResponse[[]SkillVersion]{
			Success: true,
			Data: []SkillVersion{
				{ID: "ver-1", SkillID: "skill-123", Version: 1, Semver: "1.0.0"},
				{ID: "ver-2", SkillID: "skill-123", Version: 2, Semver: "1.1.0"},
			},
		})
	}))
	defer server.Close()

	c := NewClient(server.URL, "lsk_test123")
	versions, err := c.ListSkillVersions(context.Background(), "skill-123")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(versions) != 2 {
		t.Fatalf("expected 2 versions, got %d", len(versions))
	}
	if versions[0].Semver != "1.0.0" {
		t.Errorf("expected semver '1.0.0', got '%s'", versions[0].Semver)
	}
}

func TestGetSkillContent(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("expected GET, got %s", r.Method)
		}
		if !hasPrefix(r.URL.Path, "/api/skills/skill-123/content") {
			t.Errorf("expected path /api/skills/skill-123/content, got %s", r.URL.Path)
		}
		if r.URL.Query().Get("version") != "2" {
			t.Errorf("expected version=2, got %s", r.URL.Query().Get("version"))
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(ApiResponse[SkillContent]{
			Success: true,
			Data: SkillContent{
				Content: "# Skill content v2",
				Format:  "markdown",
				Version: 2,
				Semver:  "1.1.0",
			},
		})
	}))
	defer server.Close()

	c := NewClient(server.URL, "lsk_test123")
	content, err := c.GetSkillContent(context.Background(), "skill-123", map[string]string{
		"version": "2",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if content.Content != "# Skill content v2" {
		t.Errorf("expected content '# Skill content v2', got '%s'", content.Content)
	}
	if content.Version != 2 {
		t.Errorf("expected version 2, got %d", content.Version)
	}
}

func TestGetSkillContent_NoParams(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.RawQuery != "" {
			t.Errorf("expected no query params, got %s", r.URL.RawQuery)
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(ApiResponse[SkillContent]{
			Success: true,
			Data: SkillContent{
				Content: "# Latest content",
				Format:  "markdown",
				Version: 1,
			},
		})
	}))
	defer server.Close()

	c := NewClient(server.URL, "lsk_test123")
	content, err := c.GetSkillContent(context.Background(), "skill-123", nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if content.Content != "# Latest content" {
		t.Errorf("unexpected content: %s", content.Content)
	}
}

func TestGetSkillAnalytics(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("expected GET, got %s", r.Method)
		}
		if r.URL.Path != "/api/skills/skill-123/analytics" {
			t.Errorf("expected /api/skills/skill-123/analytics, got %s", r.URL.Path)
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(ApiResponse[SkillAnalytics]{
			Success: true,
			Data: SkillAnalytics{
				TotalDownloads: 150,
				UniqueUsers:    42,
				UniqueIPs:      35,
			},
		})
	}))
	defer server.Close()

	c := NewClient(server.URL, "lsk_test123")
	analytics, err := c.GetSkillAnalytics(context.Background(), "skill-123")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if analytics.TotalDownloads != 150 {
		t.Errorf("expected 150 downloads, got %d", analytics.TotalDownloads)
	}
	if analytics.UniqueUsers != 42 {
		t.Errorf("expected 42 unique users, got %d", analytics.UniqueUsers)
	}
}

func TestRevertSkill(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("expected POST, got %s", r.Method)
		}
		if r.URL.Path != "/api/skills/skill-123/revert" {
			t.Errorf("expected /api/skills/skill-123/revert, got %s", r.URL.Path)
		}

		var body map[string]int
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			t.Fatalf("failed to decode body: %v", err)
		}
		if body["version"] != 1 {
			t.Errorf("expected version 1, got %d", body["version"])
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(ApiResponse[Skill]{
			Success: true,
			Data: Skill{
				ID:             "skill-123",
				Name:           "test-skill",
				CurrentVersion: 1,
			},
		})
	}))
	defer server.Close()

	c := NewClient(server.URL, "lsk_test123")
	skill, err := c.RevertSkill(context.Background(), "skill-123", 1)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if skill.CurrentVersion != 1 {
		t.Errorf("expected version 1, got %d", skill.CurrentVersion)
	}
}

func TestGetSkillManifest(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("expected GET, got %s", r.Method)
		}
		if r.URL.Path != "/api/skills/skill-123/manifest" {
			t.Errorf("expected /api/skills/skill-123/manifest, got %s", r.URL.Path)
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(ApiResponse[PackageManifest]{
			Success: true,
			Data: PackageManifest{
				Name:        "test-skill",
				Description: "A test skill",
				Version:     "1.0.0",
				Files:       []string{"README.md", "skill.md"},
			},
		})
	}))
	defer server.Close()

	c := NewClient(server.URL, "lsk_test123")
	manifest, err := c.GetSkillManifest(context.Background(), "skill-123")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if manifest.Name != "test-skill" {
		t.Errorf("expected name 'test-skill', got '%s'", manifest.Name)
	}
	if len(manifest.Files) != 2 {
		t.Errorf("expected 2 files, got %d", len(manifest.Files))
	}
}

func TestExplore(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("expected GET, got %s", r.Method)
		}
		if !hasPrefix(r.URL.Path, "/api/explore") {
			t.Errorf("expected path /api/explore, got %s", r.URL.Path)
		}
		if r.URL.Query().Get("query") != "terraform" {
			t.Errorf("expected query=terraform, got %s", r.URL.Query().Get("query"))
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(ApiResponse[[]ExploreSkill]{
			Success: true,
			Data: []ExploreSkill{
				{
					ID:             "skill-1",
					Name:           "terraform-helper",
					Slug:           "terraform-helper",
					AuthorName:     "testuser",
					AuthorUsername: "testuser",
					Downloads:      100,
				},
			},
		})
	}))
	defer server.Close()

	c := NewClient(server.URL, "lsk_test123")
	skills, err := c.Explore(context.Background(), map[string]string{
		"query": "terraform",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(skills) != 1 {
		t.Fatalf("expected 1 skill, got %d", len(skills))
	}
	if skills[0].Name != "terraform-helper" {
		t.Errorf("expected name 'terraform-helper', got '%s'", skills[0].Name)
	}
	if skills[0].Downloads != 100 {
		t.Errorf("expected 100 downloads, got %d", skills[0].Downloads)
	}
}

func TestExplore_NoParams(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.RawQuery != "" {
			t.Errorf("expected no query params, got %s", r.URL.RawQuery)
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(ApiResponse[[]ExploreSkill]{
			Success: true,
			Data:    []ExploreSkill{},
		})
	}))
	defer server.Close()

	c := NewClient(server.URL, "lsk_test123")
	skills, err := c.Explore(context.Background(), nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(skills) != 0 {
		t.Errorf("expected 0 skills, got %d", len(skills))
	}
}
