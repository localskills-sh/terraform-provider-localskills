package skill_version

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/localskills/terraform-provider-localskills/internal/client"
)

var (
	_ resource.Resource                = &SkillVersionResource{}
	_ resource.ResourceWithConfigure   = &SkillVersionResource{}
	_ resource.ResourceWithImportState = &SkillVersionResource{}
)

type SkillVersionResource struct {
	client *client.Client
}

func NewResource() resource.Resource {
	return &SkillVersionResource{}
}

func (r *SkillVersionResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_skill_version"
}

func (r *SkillVersionResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Publishes a new version of a skill on localskills.sh. Skill versions are **immutable** — all user-settable fields trigger replacement on change. Use `bump` to auto-increment the semver (`major`, `minor`, or `patch`) or set `semver` explicitly.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: "The ID of the skill version.",
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"skill_id": schema.StringAttribute{
				Description: "The ID of the skill this version belongs to.",
				Required:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"content": schema.StringAttribute{
				Description: "The content of this version.",
				Required:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"message": schema.StringAttribute{
				Description: "A message describing this version.",
				Optional:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"semver": schema.StringAttribute{
				Description: "The semantic version string (e.g., '1.2.0').",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"bump": schema.StringAttribute{
				Description: "The semver bump type: 'major', 'minor', or 'patch'.",
				Optional:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Validators: []validator.String{
					stringvalidator.OneOf("major", "minor", "patch"),
				},
			},
			"version": schema.Int64Attribute{
				Description: "The version number.",
				Computed:    true,
			},
			"content_hash": schema.StringAttribute{
				Description: "The hash of the content.",
				Computed:    true,
			},
			"format": schema.StringAttribute{
				Description: "The format of the content.",
				Computed:    true,
			},
			"file_count": schema.Int64Attribute{
				Description: "The number of files in this version.",
				Computed:    true,
			},
			"created_by": schema.StringAttribute{
				Description: "The user ID who created this version.",
				Computed:    true,
			},
			"created_at": schema.StringAttribute{
				Description: "The timestamp when this version was created.",
				Computed:    true,
			},
		},
	}
}

func (r *SkillVersionResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}
	c, ok := req.ProviderData.(*client.Client)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Resource Configure Type",
			fmt.Sprintf("Expected *client.Client, got: %T", req.ProviderData),
		)
		return
	}
	r.client = c
}

func (r *SkillVersionResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan SkillVersionModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	createReq := client.CreateSkillVersionRequest{
		Content: plan.Content.ValueString(),
	}
	if !plan.Message.IsNull() && !plan.Message.IsUnknown() {
		createReq.Message = plan.Message.ValueString()
	}
	if !plan.Semver.IsNull() && !plan.Semver.IsUnknown() {
		createReq.Semver = plan.Semver.ValueString()
	}
	if !plan.Bump.IsNull() && !plan.Bump.IsUnknown() {
		createReq.Bump = plan.Bump.ValueString()
	}

	ver, err := r.client.CreateSkillVersion(ctx, plan.SkillID.ValueString(), createReq)
	if err != nil {
		resp.Diagnostics.AddError("Error creating skill version", err.Error())
		return
	}

	plan.ID = types.StringValue(ver.ID)
	plan.SkillID = types.StringValue(ver.SkillID)
	plan.Version = types.Int64Value(int64(ver.Version))
	plan.Semver = types.StringValue(ver.Semver)
	plan.ContentHash = types.StringValue(ver.ContentHash)
	plan.Format = types.StringValue(ver.Format)
	plan.FileCount = types.Int64Value(int64(ver.FileCount))
	plan.CreatedBy = types.StringValue(ver.CreatedBy)
	plan.CreatedAt = types.StringValue(ver.CreatedAt)

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *SkillVersionResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state SkillVersionModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	versions, err := r.client.ListSkillVersions(ctx, state.SkillID.ValueString())
	if err != nil {
		if client.IsNotFound(err) {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError("Error reading skill versions", err.Error())
		return
	}

	versionNum := state.Version.ValueInt64()
	var found *client.SkillVersion
	for i := range versions {
		if int64(versions[i].Version) == versionNum {
			found = &versions[i]
			break
		}
	}

	if found == nil {
		resp.State.RemoveResource(ctx)
		return
	}

	preservedContent := state.Content
	preservedMessage := state.Message
	preservedBump := state.Bump

	state.ID = types.StringValue(found.ID)
	state.SkillID = types.StringValue(found.SkillID)
	state.Version = types.Int64Value(int64(found.Version))
	state.Semver = types.StringValue(found.Semver)
	state.ContentHash = types.StringValue(found.ContentHash)
	state.Format = types.StringValue(found.Format)
	state.FileCount = types.Int64Value(int64(found.FileCount))
	state.CreatedBy = types.StringValue(found.CreatedBy)
	state.CreatedAt = types.StringValue(found.CreatedAt)
	state.Content = preservedContent
	state.Message = preservedMessage
	state.Bump = preservedBump

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *SkillVersionResource) Update(_ context.Context, _ resource.UpdateRequest, resp *resource.UpdateResponse) {
	resp.Diagnostics.AddError(
		"Update Not Supported",
		"Skill versions are immutable. All changes require replacement.",
	)
}

func (r *SkillVersionResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state SkillVersionModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	skill, err := r.client.GetSkill(ctx, state.SkillID.ValueString())
	if err != nil {
		if client.IsNotFound(err) {
			return
		}
		resp.Diagnostics.AddError("Error reading skill for version delete", err.Error())
		return
	}

	versionNum := state.Version.ValueInt64()
	if int64(skill.CurrentVersion) == versionNum && versionNum > 1 {
		_, err = r.client.RevertSkill(ctx, state.SkillID.ValueString(), int(versionNum-1))
		if err != nil {
			resp.Diagnostics.AddError("Error reverting skill version", err.Error())
			return
		}
	}
}

func (r *SkillVersionResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	parts := strings.SplitN(req.ID, "/", 2)
	if len(parts) != 2 {
		resp.Diagnostics.AddError(
			"Invalid Import ID",
			fmt.Sprintf("Expected format: skill_id/version_number, got: %s", req.ID),
		)
		return
	}

	skillID := parts[0]
	versionNum, err := strconv.ParseInt(parts[1], 10, 64)
	if err != nil {
		resp.Diagnostics.AddError(
			"Invalid Version Number",
			fmt.Sprintf("Expected integer version number, got: %s", parts[1]),
		)
		return
	}

	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("skill_id"), types.StringValue(skillID))...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("version"), types.Int64Value(versionNum))...)
}
