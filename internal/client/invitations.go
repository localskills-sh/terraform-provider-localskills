package client

import (
	"context"
	"fmt"
	"net/http"
)

func (c *Client) CreateInvitation(ctx context.Context, tenantID string, req CreateInvitationRequest) (*TenantInvitation, error) {
	path := fmt.Sprintf("/api/tenants/%s/invitations", tenantID)
	return DoJSON[TenantInvitation](c, ctx, http.MethodPost, path, req)
}

func (c *Client) ListInvitations(ctx context.Context, tenantID string) ([]TenantInvitation, error) {
	path := fmt.Sprintf("/api/tenants/%s/invitations", tenantID)
	result, err := DoJSON[[]TenantInvitation](c, ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, err
	}
	return *result, nil
}
