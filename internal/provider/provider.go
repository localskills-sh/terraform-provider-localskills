package provider

import (
	"context"
	"os"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/localskills-sh/terraform-provider/internal/client"

	// Resources
	oidctrustpolicyresource "github.com/localskills-sh/terraform-provider/internal/resources/oidc_trust_policy"
	scimtokenresource "github.com/localskills-sh/terraform-provider/internal/resources/scim_token"
	skillresource "github.com/localskills-sh/terraform-provider/internal/resources/skill"
	skillversionresource "github.com/localskills-sh/terraform-provider/internal/resources/skill_version"
	ssoconnectionresource "github.com/localskills-sh/terraform-provider/internal/resources/sso_connection"
	teamresource "github.com/localskills-sh/terraform-provider/internal/resources/team"
	teaminvitationresource "github.com/localskills-sh/terraform-provider/internal/resources/team_invitation"
	teamtokenresource "github.com/localskills-sh/terraform-provider/internal/resources/team_token"
	usertokenresource "github.com/localskills-sh/terraform-provider/internal/resources/user_token"

	// Data Sources
	exploreds "github.com/localskills-sh/terraform-provider/internal/datasources/explore"
	oidctrustpoliciesds "github.com/localskills-sh/terraform-provider/internal/datasources/oidc_trust_policies"
	scimtokensds "github.com/localskills-sh/terraform-provider/internal/datasources/scim_tokens"
	skillds "github.com/localskills-sh/terraform-provider/internal/datasources/skill"
	skillanalyticsds "github.com/localskills-sh/terraform-provider/internal/datasources/skill_analytics"
	skillcontentds "github.com/localskills-sh/terraform-provider/internal/datasources/skill_content"
	skillmanifestds "github.com/localskills-sh/terraform-provider/internal/datasources/skill_manifest"
	skillversionsds "github.com/localskills-sh/terraform-provider/internal/datasources/skill_versions"
	skillsds "github.com/localskills-sh/terraform-provider/internal/datasources/skills"
	ssoconnectionds "github.com/localskills-sh/terraform-provider/internal/datasources/sso_connection"
	teamds "github.com/localskills-sh/terraform-provider/internal/datasources/team"
	teamauditlogds "github.com/localskills-sh/terraform-provider/internal/datasources/team_audit_log"
	teaminvitationsds "github.com/localskills-sh/terraform-provider/internal/datasources/team_invitations"
	teamtokensds "github.com/localskills-sh/terraform-provider/internal/datasources/team_tokens"
	teamsds "github.com/localskills-sh/terraform-provider/internal/datasources/teams"
	userauditlogds "github.com/localskills-sh/terraform-provider/internal/datasources/user_audit_log"
	userprofileds "github.com/localskills-sh/terraform-provider/internal/datasources/user_profile"
	usertokensds "github.com/localskills-sh/terraform-provider/internal/datasources/user_tokens"
)

var _ provider.Provider = &LocalskillsProvider{}

type LocalskillsProvider struct {
	version string
}

type LocalskillsProviderModel struct {
	BaseURL  types.String `tfsdk:"base_url"`
	ApiToken types.String `tfsdk:"api_token"`
}

func (p *LocalskillsProvider) Metadata(_ context.Context, _ provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "localskills"
	resp.Version = p.version
}

func (p *LocalskillsProvider) Schema(_ context.Context, _ provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "The Localskills provider manages resources on the [localskills.sh](https://localskills.sh) skill sharing platform. It supports managing skills, teams, API tokens, OIDC trust policies, SAML SSO connections, and SCIM provisioning tokens.",
		Attributes: map[string]schema.Attribute{
			"base_url": schema.StringAttribute{
				Description: "The base URL of the Localskills API. Defaults to https://localskills.sh. Can also be set with the LOCALSKILLS_BASE_URL environment variable.",
				Optional:    true,
			},
			"api_token": schema.StringAttribute{
				Description: "The API token for authenticating with the Localskills API. Must start with 'lsk_'. Can also be set with the LOCALSKILLS_API_TOKEN environment variable.",
				Required:    true,
				Sensitive:   true,
			},
		},
	}
}

func (p *LocalskillsProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
	var config LocalskillsProviderModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if config.BaseURL.IsUnknown() {
		resp.Diagnostics.AddWarning(
			"Unknown Localskills Base URL",
			"The provider cannot create the Localskills API client as there is an unknown configuration value for the base_url. "+
				"Set the value statically in the configuration, or use the LOCALSKILLS_BASE_URL environment variable.",
		)
	}
	if config.ApiToken.IsUnknown() {
		resp.Diagnostics.AddWarning(
			"Unknown Localskills API Token",
			"The provider cannot create the Localskills API client as there is an unknown configuration value for the api_token. "+
				"Set the value statically in the configuration, or use the LOCALSKILLS_API_TOKEN environment variable.",
		)
	}
	if resp.Diagnostics.HasError() {
		return
	}

	baseURL := "https://localskills.sh"
	if envURL := os.Getenv("LOCALSKILLS_BASE_URL"); envURL != "" {
		baseURL = envURL
	}
	if !config.BaseURL.IsNull() {
		baseURL = config.BaseURL.ValueString()
	}

	apiToken := os.Getenv("LOCALSKILLS_API_TOKEN")
	if !config.ApiToken.IsNull() {
		apiToken = config.ApiToken.ValueString()
	}

	if apiToken == "" {
		resp.Diagnostics.AddError(
			"Missing API Token",
			"The provider requires an API token to authenticate with the Localskills API. "+
				"Set the api_token attribute in the provider configuration or use the LOCALSKILLS_API_TOKEN environment variable.",
		)
		return
	}

	if !strings.HasPrefix(apiToken, "lsk_") {
		resp.Diagnostics.AddError(
			"Invalid API Token",
			"The API token must start with 'lsk_'. Please provide a valid Localskills API token.",
		)
		return
	}

	tflog.Debug(ctx, "Creating Localskills client", map[string]interface{}{
		"base_url": baseURL,
	})

	c := client.NewClient(baseURL, apiToken)
	resp.ResourceData = c
	resp.DataSourceData = c
}

func (p *LocalskillsProvider) Resources(_ context.Context) []func() resource.Resource {
	return []func() resource.Resource{
		skillresource.NewResource,
		skillversionresource.NewResource,
		teamresource.NewResource,
		teaminvitationresource.NewResource,
		teamtokenresource.NewResource,
		usertokenresource.NewResource,
		oidctrustpolicyresource.NewResource,
		ssoconnectionresource.NewResource,
		scimtokenresource.NewResource,
	}
}

func (p *LocalskillsProvider) DataSources(_ context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{
		skillds.NewDataSource,
		skillsds.NewDataSource,
		skillversionsds.NewDataSource,
		skillcontentds.NewDataSource,
		skillanalyticsds.NewDataSource,
		skillmanifestds.NewDataSource,
		exploreds.NewDataSource,
		teamsds.NewDataSource,
		teamds.NewDataSource,
		teaminvitationsds.NewDataSource,
		usertokensds.NewDataSource,
		teamtokensds.NewDataSource,
		oidctrustpoliciesds.NewDataSource,
		ssoconnectionds.NewDataSource,
		scimtokensds.NewDataSource,
		userprofileds.NewDataSource,
		userauditlogds.NewDataSource,
		teamauditlogds.NewDataSource,
	}
}

func New(version string) func() provider.Provider {
	return func() provider.Provider {
		return &LocalskillsProvider{
			version: version,
		}
	}
}
