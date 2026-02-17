package skill

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/localskills/terraform-provider-localskills/internal/client"
)

var (
	_ resource.Resource                = &SkillResource{}
	_ resource.ResourceWithConfigure   = &SkillResource{}
	_ resource.ResourceWithImportState = &SkillResource{}
)

type SkillResource struct {
	client *client.Client
}

func NewResource() resource.Resource {
	return &SkillResource{}
}

func (r *SkillResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_skill"
}

func (r *SkillResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Manages a skill on localskills.sh. Skills can be of type `skill` (reusable code/prompts) or `rule` (configuration rules). The `content` attribute is set at creation time; use `localskills_skill_version` to publish subsequent versions.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: "The internal ID of the skill.",
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"public_id": schema.StringAttribute{
				Description: "The public ID of the skill.",
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"tenant_id": schema.StringAttribute{
				Description: "The tenant (team) ID that owns this skill.",
				Required:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"name": schema.StringAttribute{
				Description: "The name of the skill.",
				Required:    true,
			},
			"slug": schema.StringAttribute{
				Description: "The URL slug of the skill.",
				Computed:    true,
			},
			"description": schema.StringAttribute{
				Description: "The description of the skill.",
				Optional:    true,
				Computed:    true,
				Default:     stringdefault.StaticString(""),
			},
			"type": schema.StringAttribute{
				Description: "The type of the skill. Must be 'skill' or 'rule'.",
				Required:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Validators: []validator.String{
					stringvalidator.OneOf("skill", "rule"),
				},
			},
			"visibility": schema.StringAttribute{
				Description: "The visibility of the skill. Must be 'public', 'private', or 'unlisted'.",
				Required:    true,
				Validators: []validator.String{
					stringvalidator.OneOf("public", "private", "unlisted"),
				},
			},
			"content": schema.StringAttribute{
				Description: "The content of the skill.",
				Required:    true,
			},
			"tags": schema.ListAttribute{
				Description: "Tags associated with the skill.",
				Optional:    true,
				ElementType: types.StringType,
			},
			"current_version": schema.Int64Attribute{
				Description: "The current version number of the skill.",
				Computed:    true,
			},
			"current_semver": schema.StringAttribute{
				Description: "The current semantic version of the skill.",
				Computed:    true,
			},
			"created_by": schema.StringAttribute{
				Description: "The user ID who created the skill.",
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"created_at": schema.StringAttribute{
				Description: "The timestamp when the skill was created.",
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"updated_at": schema.StringAttribute{
				Description: "The timestamp when the skill was last updated.",
				Computed:    true,
			},
		},
	}
}

func (r *SkillResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *SkillResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan SkillModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var tags []string
	if !plan.Tags.IsNull() && !plan.Tags.IsUnknown() {
		resp.Diagnostics.Append(plan.Tags.ElementsAs(ctx, &tags, false)...)
		if resp.Diagnostics.HasError() {
			return
		}
	}

	createReq := client.CreateSkillRequest{
		Name:       plan.Name.ValueString(),
		Type:       plan.Type.ValueString(),
		Visibility: plan.Visibility.ValueString(),
		Content:    plan.Content.ValueString(),
		TenantID:   plan.TenantID.ValueString(),
		Tags:       tags,
	}
	if !plan.Description.IsNull() && !plan.Description.IsUnknown() {
		createReq.Description = plan.Description.ValueString()
	}

	skill, err := r.client.CreateSkill(ctx, createReq)
	if err != nil {
		resp.Diagnostics.AddError("Error creating skill", err.Error())
		return
	}

	mapSkillToState(ctx, &plan, skill, &resp.Diagnostics)
	plan.Content = types.StringValue(createReq.Content)

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *SkillResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state SkillModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	skill, err := r.client.GetSkill(ctx, state.ID.ValueString())
	if err != nil {
		if client.IsNotFound(err) {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError("Error reading skill", err.Error())
		return
	}

	preservedContent := state.Content
	mapSkillWithVersionToState(ctx, &state, skill, &resp.Diagnostics)
	state.Content = preservedContent

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *SkillResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan SkillModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var state SkillModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	name := plan.Name.ValueString()
	description := plan.Description.ValueString()
	visibility := plan.Visibility.ValueString()

	var tags []string
	if !plan.Tags.IsNull() && !plan.Tags.IsUnknown() {
		resp.Diagnostics.Append(plan.Tags.ElementsAs(ctx, &tags, false)...)
		if resp.Diagnostics.HasError() {
			return
		}
	}

	updateReq := client.UpdateSkillRequest{
		Name:        &name,
		Description: &description,
		Visibility:  &visibility,
		Tags:        tags,
	}

	skill, err := r.client.UpdateSkill(ctx, state.ID.ValueString(), updateReq)
	if err != nil {
		resp.Diagnostics.AddError("Error updating skill", err.Error())
		return
	}

	mapSkillToState(ctx, &plan, skill, &resp.Diagnostics)
	plan.Content = types.StringValue(plan.Content.ValueString())

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *SkillResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state SkillModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	err := r.client.DeleteSkill(ctx, state.ID.ValueString())
	if err != nil {
		if client.IsNotFound(err) {
			return
		}
		resp.Diagnostics.AddError("Error deleting skill", err.Error())
	}
}

func (r *SkillResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

func mapSkillToState(ctx context.Context, state *SkillModel, skill *client.Skill, diags *diag.Diagnostics) {
	state.ID = types.StringValue(skill.ID)
	state.PublicID = types.StringValue(skill.PublicID)
	state.TenantID = types.StringValue(skill.TenantID)
	state.Name = types.StringValue(skill.Name)
	state.Slug = types.StringValue(skill.Slug)
	state.Description = types.StringValue(skill.Description)
	state.Type = types.StringValue(skill.Type)
	state.Visibility = types.StringValue(skill.Visibility)
	state.CurrentVersion = types.Int64Value(int64(skill.CurrentVersion))
	state.CurrentSemver = types.StringValue(skill.CurrentSemver)
	state.CreatedBy = types.StringValue(skill.CreatedBy)
	state.CreatedAt = types.StringValue(skill.CreatedAt)
	state.UpdatedAt = types.StringValue(skill.UpdatedAt)

	if skill.Tags != nil {
		tagValues, d := types.ListValueFrom(ctx, types.StringType, skill.Tags)
		diags.Append(d...)
		state.Tags = tagValues
	} else {
		state.Tags = types.ListNull(types.StringType)
	}
}

func mapSkillWithVersionToState(ctx context.Context, state *SkillModel, skill *client.SkillWithVersion, diags *diag.Diagnostics) {
	mapSkillToState(ctx, state, &skill.Skill, diags)
}
