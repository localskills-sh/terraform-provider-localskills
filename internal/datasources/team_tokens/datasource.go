package team_tokens

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/localskills/terraform-provider-localskills/internal/client"
)

var _ datasource.DataSource = &teamTokensDataSource{}

type teamTokensDataSource struct {
	client *client.Client
}

func NewDataSource() datasource.DataSource {
	return &teamTokensDataSource{}
}

func (d *teamTokensDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_team_tokens"
}

func (d *teamTokensDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Lists all API tokens for a team.",
		Attributes: map[string]schema.Attribute{
			"tenant_id": schema.StringAttribute{
				Description: "The tenant (team) ID.",
				Required:    true,
			},
			"tokens": schema.ListNestedAttribute{
				Description: "List of team API tokens.",
				Computed:    true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.StringAttribute{
							Description: "The unique identifier of the token.",
							Computed:    true,
						},
						"name": schema.StringAttribute{
							Description: "The name of the token.",
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
						"created_by_name": schema.StringAttribute{
							Description: "Name of the user who created the token.",
							Computed:    true,
						},
						"created_by_email": schema.StringAttribute{
							Description: "Email of the user who created the token.",
							Computed:    true,
						},
					},
				},
			},
		},
	}
}

func (d *teamTokensDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *teamTokensDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var config TeamTokensModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)
	if resp.Diagnostics.HasError() {
		return
	}

	tokens, err := d.client.ListTeamTokens(ctx, config.TenantID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Error reading team tokens", err.Error())
		return
	}

	var state TeamTokensModel
	state.TenantID = config.TenantID
	state.Tokens = make([]TeamTokenItemModel, len(tokens))
	for i, t := range tokens {
		state.Tokens[i] = TeamTokenItemModel{
			ID:             types.StringValue(t.ID),
			Name:           types.StringValue(t.Name),
			CreatedAt:      types.StringValue(t.CreatedAt),
			CreatedByEmail: types.StringValue(t.CreatedByEmail),
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
		if t.CreatedByName != nil {
			state.Tokens[i].CreatedByName = types.StringValue(*t.CreatedByName)
		} else {
			state.Tokens[i].CreatedByName = types.StringNull()
		}
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}
