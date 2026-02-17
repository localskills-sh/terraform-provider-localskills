package team

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/localskills/terraform-provider-localskills/internal/client"
)

var (
	_ resource.Resource                = &TeamResource{}
	_ resource.ResourceWithConfigure   = &TeamResource{}
	_ resource.ResourceWithImportState = &TeamResource{}
)

type TeamResource struct {
	client *client.Client
}

func NewResource() resource.Resource {
	return &TeamResource{}
}

func (r *TeamResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_team"
}

func (r *TeamResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Manages a team (tenant) on localskills.sh.\n\n~> **Note:** Deleting this resource only removes it from Terraform state. The team is **not** deleted from localskills.sh as the API does not support team deletion.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: "The unique identifier of the team.",
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"name": schema.StringAttribute{
				Description: "The name of the team.",
				Required:    true,
			},
			"slug": schema.StringAttribute{
				Description: "The URL-friendly slug of the team.",
				Optional:    true,
				Computed:    true,
			},
			"description": schema.StringAttribute{
				Description: "A description of the team.",
				Optional:    true,
				Computed:    true,
			},
			"created_at": schema.StringAttribute{
				Description: "The timestamp when the team was created.",
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"updated_at": schema.StringAttribute{
				Description: "The timestamp when the team was last updated.",
				Computed:    true,
			},
		},
	}
}

func (r *TeamResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *TeamResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan TeamModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	tenant, err := r.client.CreateTenant(ctx, client.CreateTenantRequest{
		Name: plan.Name.ValueString(),
	})
	if err != nil {
		resp.Diagnostics.AddError("Error creating team", err.Error())
		return
	}

	// If slug or description were specified, update them
	if !plan.Slug.IsNull() || !plan.Description.IsNull() {
		updateReq := client.UpdateTenantRequest{}
		if !plan.Slug.IsNull() {
			slug := plan.Slug.ValueString()
			updateReq.Slug = &slug
		}
		if !plan.Description.IsNull() {
			desc := plan.Description.ValueString()
			updateReq.Description = &desc
		}

		updated, err := r.client.UpdateTenant(ctx, tenant.ID, updateReq)
		if err != nil {
			resp.Diagnostics.AddError("Error updating team after creation", err.Error())
			return
		}
		tenant = updated
	}

	plan.ID = types.StringValue(tenant.ID)
	plan.Name = types.StringValue(tenant.Name)
	plan.Slug = types.StringValue(tenant.Slug)
	plan.Description = types.StringValue(tenant.Description)
	plan.CreatedAt = types.StringValue(tenant.CreatedAt)
	plan.UpdatedAt = types.StringValue(tenant.UpdatedAt)

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *TeamResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state TeamModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	tenants, err := r.client.ListTenants(ctx)
	if err != nil {
		resp.Diagnostics.AddError("Error reading team", err.Error())
		return
	}

	var found bool
	for _, t := range tenants {
		if t.ID == state.ID.ValueString() {
			state.Name = types.StringValue(t.Name)
			state.Slug = types.StringValue(t.Slug)
			state.Description = types.StringValue(t.Description)
			state.CreatedAt = types.StringValue(t.CreatedAt)
			state.UpdatedAt = types.StringValue(t.UpdatedAt)
			found = true
			break
		}
	}

	if !found {
		tflog.Warn(ctx, "Team not found, removing from state", map[string]interface{}{
			"id": state.ID.ValueString(),
		})
		resp.State.RemoveResource(ctx)
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *TeamResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan TeamModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var state TeamModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	updateReq := client.UpdateTenantRequest{}

	if !plan.Name.Equal(state.Name) {
		name := plan.Name.ValueString()
		updateReq.Name = &name
	}
	if !plan.Slug.Equal(state.Slug) {
		slug := plan.Slug.ValueString()
		updateReq.Slug = &slug
	}
	if !plan.Description.Equal(state.Description) {
		desc := plan.Description.ValueString()
		updateReq.Description = &desc
	}

	tenant, err := r.client.UpdateTenant(ctx, state.ID.ValueString(), updateReq)
	if err != nil {
		resp.Diagnostics.AddError("Error updating team", err.Error())
		return
	}

	plan.ID = state.ID
	plan.Name = types.StringValue(tenant.Name)
	plan.Slug = types.StringValue(tenant.Slug)
	plan.Description = types.StringValue(tenant.Description)
	plan.CreatedAt = types.StringValue(tenant.CreatedAt)
	plan.UpdatedAt = types.StringValue(tenant.UpdatedAt)

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *TeamResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state TeamModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Warn(ctx, "Team deletion is not supported by the API. Removing team from Terraform state only.", map[string]interface{}{
		"id":   state.ID.ValueString(),
		"name": state.Name.ValueString(),
	})
}

func (r *TeamResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
