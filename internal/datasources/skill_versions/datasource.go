package skill_versions

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/localskills/terraform-provider-localskills/internal/client"
)

var _ datasource.DataSource = &SkillVersionsDataSource{}
var _ datasource.DataSourceWithConfigure = &SkillVersionsDataSource{}

type SkillVersionsDataSource struct {
	client *client.Client
}

func NewDataSource() datasource.DataSource {
	return &SkillVersionsDataSource{}
}

func (d *SkillVersionsDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_skill_versions"
}

func (d *SkillVersionsDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Lists all versions of a skill.",
		Attributes: map[string]schema.Attribute{
			"skill_id": schema.StringAttribute{
				Description: "The ID of the skill.",
				Required:    true,
			},
			"versions": schema.ListNestedAttribute{
				Description: "The list of versions.",
				Computed:    true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id":           schema.StringAttribute{Computed: true},
						"skill_id":     schema.StringAttribute{Computed: true},
						"version":      schema.Int64Attribute{Computed: true},
						"semver":       schema.StringAttribute{Computed: true},
						"content_hash": schema.StringAttribute{Computed: true},
						"message":      schema.StringAttribute{Computed: true},
						"format":       schema.StringAttribute{Computed: true},
						"file_count":   schema.Int64Attribute{Computed: true},
						"created_by":   schema.StringAttribute{Computed: true},
						"created_at":   schema.StringAttribute{Computed: true},
					},
				},
			},
		},
	}
}

func (d *SkillVersionsDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *SkillVersionsDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data SkillVersionsDataSourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	versions, err := d.client.ListSkillVersions(ctx, data.SkillID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Error listing skill versions", err.Error())
		return
	}

	data.Versions = make([]SkillVersionItemModel, len(versions))
	for i, v := range versions {
		data.Versions[i] = SkillVersionItemModel{
			ID:          types.StringValue(v.ID),
			SkillID:     types.StringValue(v.SkillID),
			Version:     types.Int64Value(int64(v.Version)),
			Semver:      types.StringValue(v.Semver),
			ContentHash: types.StringValue(v.ContentHash),
			Message:     types.StringValue(v.Message),
			Format:      types.StringValue(v.Format),
			FileCount:   types.Int64Value(int64(v.FileCount)),
			CreatedBy:   types.StringValue(v.CreatedBy),
			CreatedAt:   types.StringValue(v.CreatedAt),
		}
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
