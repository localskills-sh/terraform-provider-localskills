package client

import (
	"context"
	"fmt"
	"net/http"
)

func (c *Client) ListOIDCPolicies(ctx context.Context, tenantID string) ([]OidcTrustPolicy, error) {
	policies, err := DoJSON[[]OidcTrustPolicy](c, ctx, http.MethodGet, fmt.Sprintf("/api/tenants/%s/oidc-policies", tenantID), nil)
	if err != nil {
		return nil, fmt.Errorf("listing OIDC policies: %w", err)
	}
	return *policies, nil
}

func (c *Client) CreateOIDCPolicy(ctx context.Context, tenantID string, req CreateOidcPolicyRequest) (*OidcTrustPolicy, error) {
	policy, err := DoJSON[OidcTrustPolicy](c, ctx, http.MethodPost, fmt.Sprintf("/api/tenants/%s/oidc-policies", tenantID), req)
	if err != nil {
		return nil, fmt.Errorf("creating OIDC policy: %w", err)
	}
	return policy, nil
}

func (c *Client) UpdateOIDCPolicy(ctx context.Context, tenantID, policyID string, req UpdateOidcPolicyRequest) (*OidcTrustPolicy, error) {
	policy, err := DoJSON[OidcTrustPolicy](c, ctx, http.MethodPatch, fmt.Sprintf("/api/tenants/%s/oidc-policies/%s", tenantID, policyID), req)
	if err != nil {
		return nil, fmt.Errorf("updating OIDC policy: %w", err)
	}
	return policy, nil
}

func (c *Client) DeleteOIDCPolicy(ctx context.Context, tenantID, policyID string) error {
	_, err := DoJSON[struct{}](c, ctx, http.MethodDelete, fmt.Sprintf("/api/tenants/%s/oidc-policies/%s", tenantID, policyID), nil)
	if err != nil {
		return fmt.Errorf("deleting OIDC policy: %w", err)
	}
	return nil
}
