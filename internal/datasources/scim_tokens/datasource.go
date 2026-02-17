package scim_tokens

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/localskills/terraform-provider-localskills/internal/client"
)

var _ datasource.DataSource = &scimTokensDataSource{}

type scimTokensDataSource struct {
	client *client.Client
}

func NewDataSource() datasource.DataSource {
	return &scimTokensDataSource{}
}

func (d *scimTokensDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_scim_tokens"
}

func (d *scimTokensDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Lists all SCIM provisioning tokens for a team.",
		Attributes: map[string]schema.Attribute{
			"tenant_id": schema.StringAttribute{
				Description: "The tenant (team) ID.",
				Required:    true,
			},
			"tokens": schema.ListNestedAttribute{
				Description: "List of SCIM tokens.",
				Computed:    true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.StringAttribute{
							Description: "The unique identifier of the SCIM token.",
							Computed:    true,
						},
						"name": schema.StringAttribute{
							Description: "The name of the SCIM token.",
							Computed:    true,
						},
						"last_used_at": schema.StringAttribute{
							Description: "When the token was last used.",
							Computed:    true,
						},
						"expires_at": schema.StringAttribute{
							Description: "When the token expires.",
							Computed:    true,
						},
						"created_at": schema.StringAttribute{
							Description: "When the token was created.",
							Computed:    true,
						},
					},
				},
			},
		},
	}
}

func (d *scimTokensDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}
	c, ok := req.ProviderData.(*client.Client)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Data Source Configure Type",
			fmt.Sprintf("Expected *client.Client, got: %T", req.ProviderData),
		)
		return
	}
	d.client = c
}

func (d *scimTokensDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var config ScimTokensModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)
	if resp.Diagnostics.HasError() {
		return
	}

	tokens, err := d.client.ListSCIMTokens(ctx, config.TenantID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Error reading SCIM tokens", err.Error())
		return
	}

	var state ScimTokensModel
	state.TenantID = config.TenantID
	state.Tokens = make([]ScimTokenItemModel, len(tokens))
	for i, t := range tokens {
		state.Tokens[i] = ScimTokenItemModel{
			ID:        types.StringValue(t.ID),
			Name:      types.StringValue(t.Name),
			CreatedAt: types.StringValue(t.CreatedAt),
		}
		if t.LastUsedAt != nil {
			state.Tokens[i].LastUsedAt = types.StringValue(*t.LastUsedAt)
		} else {
			state.Tokens[i].LastUsedAt = types.StringNull()
		}
		if t.ExpiresAt != nil {
			state.Tokens[i].ExpiresAt = types.StringValue(*t.ExpiresAt)
		} else {
			state.Tokens[i].ExpiresAt = types.StringNull()
		}
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}
