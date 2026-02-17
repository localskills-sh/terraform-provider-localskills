package skill_content

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/localskills/terraform-provider-localskills/internal/client"
)

var _ datasource.DataSource = &SkillContentDataSource{}
var _ datasource.DataSourceWithConfigure = &SkillContentDataSource{}

type SkillContentDataSource struct {
	client *client.Client
}

func NewDataSource() datasource.DataSource {
	return &SkillContentDataSource{}
}

func (d *SkillContentDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_skill_content"
}

func (d *SkillContentDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Fetches the content of a skill, optionally at a specific version. Supports pinning by `version` number, exact `semver`, or a semver `range` (e.g. `~> 1.0`).",
		Attributes: map[string]schema.Attribute{
			"skill_id": schema.StringAttribute{
				Description: "The ID of the skill.",
				Required:    true,
			},
			"version": schema.StringAttribute{
				Description: "The version number to fetch.",
				Optional:    true,
			},
			"semver": schema.StringAttribute{
				Description: "The semantic version to fetch.",
				Optional:    true,
			},
			"range": schema.StringAttribute{
				Description: "The semver range to match.",
				Optional:    true,
			},
			"content": schema.StringAttribute{
				Description: "The content of the skill.",
				Computed:    true,
			},
			"format": schema.StringAttribute{
				Description: "The format of the content.",
				Computed:    true,
			},
			"skill_name": schema.StringAttribute{
				Description: "The name of the skill.",
				Computed:    true,
			},
			"skill_slug": schema.StringAttribute{
				Description: "The slug of the skill.",
				Computed:    true,
			},
			"fetch_version": schema.Int64Attribute{
				Description: "The version number that was fetched.",
				Computed:    true,
			},
			"fetch_semver": schema.StringAttribute{
				Description: "The semantic version that was fetched.",
				Computed:    true,
			},
		},
	}
}

func (d *SkillContentDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *SkillContentDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data SkillContentDataSourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	params := map[string]string{}
	if !data.Version.IsNull() && !data.Version.IsUnknown() {
		params["version"] = data.Version.ValueString()
	}
	if !data.Semver.IsNull() && !data.Semver.IsUnknown() {
		params["semver"] = data.Semver.ValueString()
	}
	if !data.Range.IsNull() && !data.Range.IsUnknown() {
		params["range"] = data.Range.ValueString()
	}

	content, err := d.client.GetSkillContent(ctx, data.SkillID.ValueString(), params)
	if err != nil {
		resp.Diagnostics.AddError("Error reading skill content", err.Error())
		return
	}

	data.Content = types.StringValue(content.Content)
	data.Format = types.StringValue(content.Format)
	data.SkillName = types.StringValue(content.Skill.Name)
	data.SkillSlug = types.StringValue(content.Skill.Slug)
	data.FetchVersion = types.Int64Value(int64(content.Version))
	data.FetchSemver = types.StringValue(content.Semver)

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
