package skill_analytics

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/localskills/terraform-provider-localskills/internal/client"
)

var _ datasource.DataSource = &SkillAnalyticsDataSource{}
var _ datasource.DataSourceWithConfigure = &SkillAnalyticsDataSource{}

type SkillAnalyticsDataSource struct {
	client *client.Client
}

func NewDataSource() datasource.DataSource {
	return &SkillAnalyticsDataSource{}
}

func (d *SkillAnalyticsDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_skill_analytics"
}

func (d *SkillAnalyticsDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Fetches download and usage analytics for a skill on localskills.sh. Returns `total_downloads`, `unique_users`, and `unique_ips`.",
		Attributes: map[string]schema.Attribute{
			"skill_id": schema.StringAttribute{
				Description: "The ID of the skill.",
				Required:    true,
			},
			"total_downloads": schema.Int64Attribute{
				Description: "Total number of downloads.",
				Computed:    true,
			},
			"unique_users": schema.Int64Attribute{
				Description: "Number of unique users.",
				Computed:    true,
			},
			"unique_ips": schema.Int64Attribute{
				Description: "Number of unique IPs.",
				Computed:    true,
			},
		},
	}
}

func (d *SkillAnalyticsDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *SkillAnalyticsDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data SkillAnalyticsDataSourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	analytics, err := d.client.GetSkillAnalytics(ctx, data.SkillID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Error reading skill analytics", err.Error())
		return
	}

	data.TotalDownloads = types.Int64Value(int64(analytics.TotalDownloads))
	data.UniqueUsers = types.Int64Value(int64(analytics.UniqueUsers))
	data.UniqueIPs = types.Int64Value(int64(analytics.UniqueIPs))

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
