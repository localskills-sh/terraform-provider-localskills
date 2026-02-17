package user_tokens

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/localskills/terraform-provider-localskills/internal/client"
)

var _ datasource.DataSource = &userTokensDataSource{}

type userTokensDataSource struct {
	client *client.Client
}

func NewDataSource() datasource.DataSource {
	return &userTokensDataSource{}
}

func (d *userTokensDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_user_tokens"
}

func (d *userTokensDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Lists all user API tokens.",
		Attributes: map[string]schema.Attribute{
			"tokens": schema.ListNestedAttribute{
				Description: "List of user API tokens.",
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
					},
				},
			},
		},
	}
}

func (d *userTokensDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *userTokensDataSource) Read(ctx context.Context, _ datasource.ReadRequest, resp *datasource.ReadResponse) {
	tokens, err := d.client.ListUserTokens(ctx)
	if err != nil {
		resp.Diagnostics.AddError("Error reading user tokens", err.Error())
		return
	}

	var state UserTokensModel
	state.Tokens = make([]UserTokenItemModel, len(tokens))
	for i, t := range tokens {
		state.Tokens[i] = UserTokenItemModel{
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
