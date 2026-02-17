package team

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/localskills/terraform-provider-localskills/internal/client"
)

var (
	_ datasource.DataSource              = &TeamDataSource{}
	_ datasource.DataSourceWithConfigure = &TeamDataSource{}
)

type TeamDataSource struct {
	client *client.Client
}

func NewDataSource() datasource.DataSource {
	return &TeamDataSource{}
}

func (d *TeamDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_team"
}

func (d *TeamDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Fetches a single team (tenant) by ID or slug.",
		Attributes: map[string]schema.Attribute{
			"team_id": schema.StringAttribute{
				Description: "The ID of the team to look up. At least one of team_id or slug must be specified.",
				Optional:    true,
			},
			"slug": schema.StringAttribute{
				Description: "The slug of the team to look up. At least one of team_id or slug must be specified.",
				Optional:    true,
				Computed:    true,
			},
			"id": schema.StringAttribute{
				Description: "The unique identifier of the team.",
				Computed:    true,
			},
			"name": schema.StringAttribute{
				Description: "The name of the team.",
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
	}
}

func (d *TeamDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *TeamDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var config TeamDataSourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if config.TeamID.IsNull() && config.Slug.IsNull() {
		resp.Diagnostics.AddError(
			"Missing Required Attribute",
			"At least one of team_id or slug must be specified.",
		)
		return
	}

	tenants, err := d.client.ListTenants(ctx)
	if err != nil {
		resp.Diagnostics.AddError("Error listing teams", err.Error())
		return
	}

	var found bool
	for _, t := range tenants {
		if (!config.TeamID.IsNull() && t.ID == config.TeamID.ValueString()) ||
			(!config.Slug.IsNull() && t.Slug == config.Slug.ValueString()) {
			config.ID = types.StringValue(t.ID)
			config.TeamID = types.StringValue(t.ID)
			config.Name = types.StringValue(t.Name)
			config.Slug = types.StringValue(t.Slug)
			config.Description = types.StringValue(t.Description)
			config.Role = types.StringValue(t.Role)
			found = true
			break
		}
	}

	if !found {
		resp.Diagnostics.AddError(
			"Team Not Found",
			"No team matching the specified criteria was found.",
		)
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &config)...)
}
