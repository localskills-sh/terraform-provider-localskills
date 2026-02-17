package client

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
)

func (c *Client) GetUserProfile(ctx context.Context) (*UserProfile, error) {
	return DoJSON[UserProfile](c, ctx, http.MethodGet, "/api/user/profile", nil)
}

func (c *Client) ListUserAuditLog(ctx context.Context, params map[string]string) (*AuditLogResponse, error) {
	path := "/api/user/audit-log"
	if len(params) > 0 {
		q := url.Values{}
		for k, v := range params {
			if v != "" {
				q.Set(k, v)
			}
		}
		path += "?" + q.Encode()
	}
	return DoJSON[AuditLogResponse](c, ctx, http.MethodGet, path, nil)
}

func (c *Client) ListTeamAuditLog(ctx context.Context, tenantID string, params map[string]string) (*AuditLogResponse, error) {
	path := fmt.Sprintf("/api/tenants/%s/audit-log", tenantID)
	if len(params) > 0 {
		q := url.Values{}
		for k, v := range params {
			if v != "" {
				q.Set(k, v)
			}
		}
		path += "?" + q.Encode()
	}
	return DoJSON[AuditLogResponse](c, ctx, http.MethodGet, path, nil)
}
