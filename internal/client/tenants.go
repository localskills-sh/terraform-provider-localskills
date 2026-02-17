package client

import (
	"context"
	"fmt"
	"net/http"
)

func (c *Client) CreateTenant(ctx context.Context, req CreateTenantRequest) (*Tenant, error) {
	return DoJSON[Tenant](c, ctx, http.MethodPost, "/api/tenants", req)
}

func (c *Client) ListTenants(ctx context.Context) ([]TenantWithRole, error) {
	result, err := DoJSON[[]TenantWithRole](c, ctx, http.MethodGet, "/api/tenants", nil)
	if err != nil {
		return nil, err
	}
	return *result, nil
}

func (c *Client) UpdateTenant(ctx context.Context, tenantID string, req UpdateTenantRequest) (*Tenant, error) {
	path := fmt.Sprintf("/api/tenants/%s", tenantID)
	return DoJSON[Tenant](c, ctx, http.MethodPatch, path, req)
}
