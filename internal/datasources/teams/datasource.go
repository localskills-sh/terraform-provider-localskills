package teams

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/localskills/terraform-provider-localskills/internal/client"
)

var (
	_ datasource.DataSource              = &TeamsDataSource{}
	_ datasource.DataSourceWithConfigure = &TeamsDataSource{}
)

type TeamsDataSource struct {
	client *client.Client
}

func NewDataSource() datasource.DataSource {
	return &TeamsDataSource{}
}

func (d *TeamsDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_teams"
}

func (d *TeamsDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Fetches all teams (tenants) accessible to the authenticated user.",
		Attributes: map[string]schema.Attribute{
			"teams": schema.ListNestedAttribute{
				Description: "List of teams.",
				Computed:    true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.StringAttribute{
							Description: "The unique identifier of the team.",
							Computed:    true,
						},
						"name": schema.StringAttribute{
							Description: "The name of the team.",
							Computed:    true,
						},
						"slug": schema.StringAttribute{
							Description: "The URL-friendly slug of the team.",
							Computed:    true,
						},
						"description": schema.StringAttribute{
							Description: "A description of the team.",
							Computed:    true,
						},
						"role": schema.StringAttribute{
							Description: "The authenticated user's role in this team.",
							Computed:    true,
						},
					},
				},
			},
		},
	}
}

func (d *TeamsDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *TeamsDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	tenants, err := d.client.ListTenants(ctx)
	if err != nil {
		resp.Diagnostics.AddError("Error listing teams", err.Error())
		return
	}

	var state TeamsDataSourceModel
	for _, t := range tenants {
		state.Teams = append(state.Teams, TeamModel{
			ID:          types.StringValue(t.ID),
			Name:        types.StringValue(t.Name),
			Slug:        types.StringValue(t.Slug),
			Description: types.StringValue(t.Description),
			Role:        types.StringValue(t.Role),
		})
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}
