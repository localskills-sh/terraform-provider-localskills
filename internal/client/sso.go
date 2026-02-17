package client

import (
	"context"
	"fmt"
	"net/http"
)

func (c *Client) GetSSOConnection(ctx context.Context, tenantID string) (*SsoConnection, error) {
	conn, err := DoJSON[SsoConnection](c, ctx, http.MethodGet, fmt.Sprintf("/api/tenants/%s/sso", tenantID), nil)
	if err != nil {
		return nil, fmt.Errorf("getting SSO connection: %w", err)
	}
	return conn, nil
}

func (c *Client) UpdateSSOConnection(ctx context.Context, tenantID string, req UpdateSsoRequest) (*SsoConnection, error) {
	conn, err := DoJSON[SsoConnection](c, ctx, http.MethodPatch, fmt.Sprintf("/api/tenants/%s/sso", tenantID), req)
	if err != nil {
		return nil, fmt.Errorf("updating SSO connection: %w", err)
	}
	return conn, nil
}
