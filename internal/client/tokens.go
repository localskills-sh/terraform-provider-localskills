package client

import (
	"context"
	"fmt"
	"net/http"
)

func (c *Client) ListUserTokens(ctx context.Context) ([]ApiToken, error) {
	tokens, err := DoJSON[[]ApiToken](c, ctx, http.MethodGet, "/api/user/tokens", nil)
	if err != nil {
		return nil, fmt.Errorf("listing user tokens: %w", err)
	}
	return *tokens, nil
}

func (c *Client) CreateUserToken(ctx context.Context, req CreateTokenRequest) (*ApiTokenWithSecret, error) {
	token, err := DoJSON[ApiTokenWithSecret](c, ctx, http.MethodPost, "/api/user/tokens", req)
	if err != nil {
		return nil, fmt.Errorf("creating user token: %w", err)
	}
	return token, nil
}

func (c *Client) DeleteUserToken(ctx context.Context, tokenID string) error {
	_, err := DoJSON[struct{}](c, ctx, http.MethodDelete, fmt.Sprintf("/api/user/tokens/%s", tokenID), nil)
	if err != nil {
		return fmt.Errorf("deleting user token: %w", err)
	}
	return nil
}

func (c *Client) ListTeamTokens(ctx context.Context, tenantID string) ([]TeamApiToken, error) {
	tokens, err := DoJSON[[]TeamApiToken](c, ctx, http.MethodGet, fmt.Sprintf("/api/tenants/%s/tokens", tenantID), nil)
	if err != nil {
		return nil, fmt.Errorf("listing team tokens: %w", err)
	}
	return *tokens, nil
}

func (c *Client) CreateTeamToken(ctx context.Context, tenantID string, req CreateTeamTokenRequest) (*TeamApiTokenWithSecret, error) {
	token, err := DoJSON[TeamApiTokenWithSecret](c, ctx, http.MethodPost, fmt.Sprintf("/api/tenants/%s/tokens", tenantID), req)
	if err != nil {
		return nil, fmt.Errorf("creating team token: %w", err)
	}
	return token, nil
}

func (c *Client) DeleteTeamToken(ctx context.Context, tenantID, tokenID string) error {
	_, err := DoJSON[struct{}](c, ctx, http.MethodDelete, fmt.Sprintf("/api/tenants/%s/tokens/%s", tenantID, tokenID), nil)
	if err != nil {
		return fmt.Errorf("deleting team token: %w", err)
	}
	return nil
}
