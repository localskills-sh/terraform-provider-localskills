package skill_manifest

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/localskills/terraform-provider-localskills/internal/client"
)

var _ datasource.DataSource = &SkillManifestDataSource{}
var _ datasource.DataSourceWithConfigure = &SkillManifestDataSource{}

type SkillManifestDataSource struct {
	client *client.Client
}

func NewDataSource() datasource.DataSource {
	return &SkillManifestDataSource{}
}

func (d *SkillManifestDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_skill_manifest"
}

func (d *SkillManifestDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Fetches the manifest of a skill package.",
		Attributes: map[string]schema.Attribute{
			"skill_id": schema.StringAttribute{
				Description: "The ID of the skill.",
				Required:    true,
			},
			"name": schema.StringAttribute{
				Description: "The name from the manifest.",
				Computed:    true,
			},
			"description": schema.StringAttribute{
				Description: "The description from the manifest.",
				Computed:    true,
			},
			"version": schema.StringAttribute{
				Description: "The version from the manifest.",
				Computed:    true,
			},
			"files": schema.ListAttribute{
				Description: "The list of files in the package.",
				Computed:    true,
				ElementType: types.StringType,
			},
		},
	}
}

func (d *SkillManifestDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *SkillManifestDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data SkillManifestDataSourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	manifest, err := d.client.GetSkillManifest(ctx, data.SkillID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Error reading skill manifest", err.Error())
		return
	}

	data.Name = types.StringValue(manifest.Name)
	data.Description = types.StringValue(manifest.Description)
	data.Version = types.StringValue(manifest.Version)

	if manifest.Files != nil {
		fileValues, diags := types.ListValueFrom(ctx, types.StringType, manifest.Files)
		resp.Diagnostics.Append(diags...)
		data.Files = fileValues
	} else {
		data.Files = types.ListNull(types.StringType)
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
