package client

import (
	"context"
	"fmt"
	"net/http"
)

func (c *Client) ListSCIMTokens(ctx context.Context, tenantID string) ([]ScimToken, error) {
	tokens, err := DoJSON[[]ScimToken](c, ctx, http.MethodGet, fmt.Sprintf("/api/tenants/%s/scim-tokens", tenantID), nil)
	if err != nil {
		return nil, fmt.Errorf("listing SCIM tokens: %w", err)
	}
	return *tokens, nil
}

func (c *Client) CreateSCIMToken(ctx context.Context, tenantID string, req CreateScimTokenRequest) (*ScimTokenWithSecret, error) {
	token, err := DoJSON[ScimTokenWithSecret](c, ctx, http.MethodPost, fmt.Sprintf("/api/tenants/%s/scim-tokens", tenantID), req)
	if err != nil {
		return nil, fmt.Errorf("creating SCIM token: %w", err)
	}
	return token, nil
}

func (c *Client) DeleteSCIMToken(ctx context.Context, tenantID, tokenID string) error {
	_, err := DoJSON[struct{}](c, ctx, http.MethodDelete, fmt.Sprintf("/api/tenants/%s/scim-tokens/%s", tenantID, tokenID), nil)
	if err != nil {
		return fmt.Errorf("deleting SCIM token: %w", err)
	}
	return nil
}
