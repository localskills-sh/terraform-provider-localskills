package skill

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/localskills/terraform-provider-localskills/internal/client"
)

var _ datasource.DataSource = &SkillDataSource{}
var _ datasource.DataSourceWithConfigure = &SkillDataSource{}

type SkillDataSource struct {
	client *client.Client
}

func NewDataSource() datasource.DataSource {
	return &SkillDataSource{}
}

func (d *SkillDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_skill"
}

func (d *SkillDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Fetches a single skill by ID.",
		Attributes: map[string]schema.Attribute{
			"skill_id": schema.StringAttribute{
				Description: "The ID of the skill to fetch.",
				Required:    true,
			},
			"id": schema.StringAttribute{
				Description: "The internal ID of the skill.",
				Computed:    true,
			},
			"public_id": schema.StringAttribute{
				Description: "The public ID of the skill.",
				Computed:    true,
			},
			"tenant_id": schema.StringAttribute{
				Description: "The tenant ID that owns this skill.",
				Computed:    true,
			},
			"name": schema.StringAttribute{
				Description: "The name of the skill.",
				Computed:    true,
			},
			"slug": schema.StringAttribute{
				Description: "The URL slug of the skill.",
				Computed:    true,
			},
			"description": schema.StringAttribute{
				Description: "The description of the skill.",
				Computed:    true,
			},
			"type": schema.StringAttribute{
				Description: "The type of the skill.",
				Computed:    true,
			},
			"visibility": schema.StringAttribute{
				Description: "The visibility of the skill.",
				Computed:    true,
			},
			"tags": schema.ListAttribute{
				Description: "Tags associated with the skill.",
				Computed:    true,
				ElementType: types.StringType,
			},
			"current_version": schema.Int64Attribute{
				Description: "The current version number.",
				Computed:    true,
			},
			"current_semver": schema.StringAttribute{
				Description: "The current semantic version.",
				Computed:    true,
			},
			"created_by": schema.StringAttribute{
				Description: "The user ID who created the skill.",
				Computed:    true,
			},
			"created_at": schema.StringAttribute{
				Description: "The creation timestamp.",
				Computed:    true,
			},
			"updated_at": schema.StringAttribute{
				Description: "The last update timestamp.",
				Computed:    true,
			},
		},
	}
}

func (d *SkillDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *SkillDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data SkillDataSourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	skill, err := d.client.GetSkill(ctx, data.SkillID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Error reading skill", err.Error())
		return
	}

	data.ID = types.StringValue(skill.ID)
	data.PublicID = types.StringValue(skill.PublicID)
	data.TenantID = types.StringValue(skill.TenantID)
	data.Name = types.StringValue(skill.Name)
	data.Slug = types.StringValue(skill.Slug)
	data.Description = types.StringValue(skill.Description)
	data.Type = types.StringValue(skill.Type)
	data.Visibility = types.StringValue(skill.Visibility)
	data.CurrentVersion = types.Int64Value(int64(skill.CurrentVersion))
	data.CurrentSemver = types.StringValue(skill.CurrentSemver)
	data.CreatedBy = types.StringValue(skill.CreatedBy)
	data.CreatedAt = types.StringValue(skill.CreatedAt)
	data.UpdatedAt = types.StringValue(skill.UpdatedAt)

	if skill.Tags != nil {
		tagValues, diags := types.ListValueFrom(ctx, types.StringType, skill.Tags)
		resp.Diagnostics.Append(diags...)
		data.Tags = tagValues
	} else {
		data.Tags = types.ListNull(types.StringType)
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
