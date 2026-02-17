package sso_connection

import (
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type SsoConnectionDataModel struct {
	TenantID     types.String `tfsdk:"tenant_id"`
	ID           types.String `tfsdk:"id"`
	DisplayName  types.String `tfsdk:"display_name"`
	MetadataURL  types.String `tfsdk:"metadata_url"`
	DefaultRole  types.String `tfsdk:"default_role"`
	EmailDomains types.List   `tfsdk:"email_domains"`
	Enabled      types.Bool   `tfsdk:"enabled"`
	RequireSso   types.Bool   `tfsdk:"require_sso"`
	IdpEntityID  types.String `tfsdk:"idp_entity_id"`
	IdpSsoURL    types.String `tfsdk:"idp_sso_url"`
	IdpSloURL    types.String `tfsdk:"idp_slo_url"`
	SpEntityID   types.String `tfsdk:"sp_entity_id"`
	SpAcsURL     types.String `tfsdk:"sp_acs_url"`
	CreatedAt    types.String `tfsdk:"created_at"`
	UpdatedAt    types.String `tfsdk:"updated_at"`
}
