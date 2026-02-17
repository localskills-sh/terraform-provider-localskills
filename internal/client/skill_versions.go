package client

import (
	"context"
	"fmt"
	"net/http"
	"strings"
)

func (c *Client) CreateSkillVersion(ctx context.Context, skillID string, req CreateSkillVersionRequest) (*SkillVersion, error) {
	return DoJSON[SkillVersion](c, ctx, http.MethodPost, fmt.Sprintf("/api/skills/%s/versions", skillID), req)
}

func (c *Client) ListSkillVersions(ctx context.Context, skillID string) ([]SkillVersion, error) {
	result, err := DoJSON[[]SkillVersion](c, ctx, http.MethodGet, fmt.Sprintf("/api/skills/%s/versions", skillID), nil)
	if err != nil {
		return nil, err
	}
	return *result, nil
}

func (c *Client) GetSkillContent(ctx context.Context, skillID string, params map[string]string) (*SkillContent, error) {
	path := fmt.Sprintf("/api/skills/%s/content", skillID)
	if len(params) > 0 {
		var parts []string
		for k, v := range params {
			parts = append(parts, fmt.Sprintf("%s=%s", k, v))
		}
		path += "?" + strings.Join(parts, "&")
	}
	return DoJSON[SkillContent](c, ctx, http.MethodGet, path, nil)
}

func (c *Client) GetSkillAnalytics(ctx context.Context, skillID string) (*SkillAnalytics, error) {
	return DoJSON[SkillAnalytics](c, ctx, http.MethodGet, fmt.Sprintf("/api/skills/%s/analytics", skillID), nil)
}

func (c *Client) RevertSkill(ctx context.Context, skillID string, version int) (*Skill, error) {
	body := map[string]int{"version": version}
	return DoJSON[Skill](c, ctx, http.MethodPost, fmt.Sprintf("/api/skills/%s/revert", skillID), body)
}

func (c *Client) GetSkillManifest(ctx context.Context, skillID string) (*PackageManifest, error) {
	return DoJSON[PackageManifest](c, ctx, http.MethodGet, fmt.Sprintf("/api/skills/%s/manifest", skillID), nil)
}

func (c *Client) Explore(ctx context.Context, params map[string]string) ([]ExploreSkill, error) {
	path := "/api/explore"
	if len(params) > 0 {
		var parts []string
		for k, v := range params {
			parts = append(parts, fmt.Sprintf("%s=%s", k, v))
		}
		path += "?" + strings.Join(parts, "&")
	}
	result, err := DoJSON[[]ExploreSkill](c, ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, err
	}
	return *result, nil
}
