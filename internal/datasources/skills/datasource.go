package skills

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/localskills/terraform-provider-localskills/internal/client"
)

var _ datasource.DataSource = &SkillsDataSource{}
var _ datasource.DataSourceWithConfigure = &SkillsDataSource{}

type SkillsDataSource struct {
	client *client.Client
}

func NewDataSource() datasource.DataSource {
	return &SkillsDataSource{}
}

func (d *SkillsDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_skills"
}

func (d *SkillsDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Lists skills with optional filters.",
		Attributes: map[string]schema.Attribute{
			"tenant_id": schema.StringAttribute{
				Description: "Filter by tenant ID.",
				Optional:    true,
			},
			"visibility": schema.StringAttribute{
				Description: "Filter by visibility.",
				Optional:    true,
			},
			"type": schema.StringAttribute{
				Description: "Filter by type.",
				Optional:    true,
			},
			"query": schema.StringAttribute{
				Description: "Search query.",
				Optional:    true,
			},
			"tag": schema.StringAttribute{
				Description: "Filter by tag.",
				Optional:    true,
			},
			"skills": schema.ListNestedAttribute{
				Description: "The list of skills.",
				Computed:    true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id":              schema.StringAttribute{Computed: true},
						"public_id":       schema.StringAttribute{Computed: true},
						"tenant_id":       schema.StringAttribute{Computed: true},
						"name":            schema.StringAttribute{Computed: true},
						"slug":            schema.StringAttribute{Computed: true},
						"description":     schema.StringAttribute{Computed: true},
						"type":            schema.StringAttribute{Computed: true},
						"visibility":      schema.StringAttribute{Computed: true},
						"tags":            schema.ListAttribute{Computed: true, ElementType: types.StringType},
						"current_version": schema.Int64Attribute{Computed: true},
						"current_semver":  schema.StringAttribute{Computed: true},
						"created_by":      schema.StringAttribute{Computed: true},
						"created_at":      schema.StringAttribute{Computed: true},
						"updated_at":      schema.StringAttribute{Computed: true},
					},
				},
			},
		},
	}
}

func (d *SkillsDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *SkillsDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data SkillsDataSourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	params := map[string]string{}
	if !data.TenantID.IsNull() && !data.TenantID.IsUnknown() {
		params["tenant_id"] = data.TenantID.ValueString()
	}
	if !data.Visibility.IsNull() && !data.Visibility.IsUnknown() {
		params["visibility"] = data.Visibility.ValueString()
	}
	if !data.Type.IsNull() && !data.Type.IsUnknown() {
		params["type"] = data.Type.ValueString()
	}
	if !data.Query.IsNull() && !data.Query.IsUnknown() {
		params["query"] = data.Query.ValueString()
	}
	if !data.Tag.IsNull() && !data.Tag.IsUnknown() {
		params["tag"] = data.Tag.ValueString()
	}

	skills, err := d.client.ListSkills(ctx, params)
	if err != nil {
		resp.Diagnostics.AddError("Error listing skills", err.Error())
		return
	}

	data.Skills = make([]SkillItemModel, len(skills))
	for i, s := range skills {
		item := SkillItemModel{
			ID:             types.StringValue(s.ID),
			PublicID:       types.StringValue(s.PublicID),
			TenantID:       types.StringValue(s.TenantID),
			Name:           types.StringValue(s.Name),
			Slug:           types.StringValue(s.Slug),
			Description:    types.StringValue(s.Description),
			Type:           types.StringValue(s.Type),
			Visibility:     types.StringValue(s.Visibility),
			CurrentVersion: types.Int64Value(int64(s.CurrentVersion)),
			CurrentSemver:  types.StringValue(s.CurrentSemver),
			CreatedBy:      types.StringValue(s.CreatedBy),
			CreatedAt:      types.StringValue(s.CreatedAt),
			UpdatedAt:      types.StringValue(s.UpdatedAt),
		}
		if s.Tags != nil {
			tagValues, diags := types.ListValueFrom(ctx, types.StringType, s.Tags)
			resp.Diagnostics.Append(diags...)
			item.Tags = tagValues
		} else {
			item.Tags = types.ListNull(types.StringType)
		}
		data.Skills[i] = item
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
