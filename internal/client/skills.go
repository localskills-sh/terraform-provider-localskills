package client

import (
	"context"
	"fmt"
	"net/http"
	"strings"
)

func (c *Client) CreateSkill(ctx context.Context, req CreateSkillRequest) (*Skill, error) {
	return DoJSON[Skill](c, ctx, http.MethodPost, "/api/skills", req)
}

func (c *Client) GetSkill(ctx context.Context, skillID string) (*SkillWithVersion, error) {
	return DoJSON[SkillWithVersion](c, ctx, http.MethodGet, fmt.Sprintf("/api/skills/%s", skillID), nil)
}

func (c *Client) UpdateSkill(ctx context.Context, skillID string, req UpdateSkillRequest) (*Skill, error) {
	return DoJSON[Skill](c, ctx, http.MethodPut, fmt.Sprintf("/api/skills/%s", skillID), req)
}

func (c *Client) DeleteSkill(ctx context.Context, skillID string) error {
	resp, err := c.doRequest(ctx, http.MethodDelete, fmt.Sprintf("/api/skills/%s", skillID), nil)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		return &ApiError{
			StatusCode: resp.StatusCode,
			Message:    fmt.Sprintf("failed to delete skill %s", skillID),
		}
	}
	return nil
}

func (c *Client) ListSkills(ctx context.Context, params map[string]string) ([]Skill, error) {
	path := "/api/skills"
	if len(params) > 0 {
		var parts []string
		for k, v := range params {
			parts = append(parts, fmt.Sprintf("%s=%s", k, v))
		}
		path += "?" + strings.Join(parts, "&")
	}
	result, err := DoJSON[[]Skill](c, ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, err
	}
	return *result, nil
}
