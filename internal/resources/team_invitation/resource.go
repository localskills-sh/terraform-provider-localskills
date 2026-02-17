package team_invitation

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/localskills/terraform-provider-localskills/internal/client"
)

var (
	_ resource.Resource                = &TeamInvitationResource{}
	_ resource.ResourceWithConfigure   = &TeamInvitationResource{}
	_ resource.ResourceWithImportState = &TeamInvitationResource{}
)

type TeamInvitationResource struct {
	client *client.Client
}

func NewResource() resource.Resource {
	return &TeamInvitationResource{}
}

func (r *TeamInvitationResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_team_invitation"
}

func (r *TeamInvitationResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Manages a team invitation on localskills.sh. Invitations are **immutable** — all changes require replacement.\n\n~> **Note:** Deleting this resource only removes it from Terraform state. The invitation is **not** revoked via the API.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: "The unique identifier of the invitation.",
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"tenant_id": schema.StringAttribute{
				Description: "The ID of the team (tenant) to invite to.",
				Required:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"email": schema.StringAttribute{
				Description: "The email address of the person to invite.",
				Required:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"role": schema.StringAttribute{
				Description: "The role to assign to the invited user. Must be one of: owner, admin, member, viewonly.",
				Required:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Validators: []validator.String{
					stringvalidator.OneOf("owner", "admin", "member", "viewonly"),
				},
			},
			"token": schema.StringAttribute{
				Description: "The invitation token.",
				Computed:    true,
				Sensitive:   true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"invited_by": schema.StringAttribute{
				Description: "The ID of the user who created the invitation.",
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"expires_at": schema.StringAttribute{
				Description: "The timestamp when the invitation expires.",
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"accepted_at": schema.StringAttribute{
				Description: "The timestamp when the invitation was accepted, if applicable.",
				Computed:    true,
			},
			"created_at": schema.StringAttribute{
				Description: "The timestamp when the invitation was created.",
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
		},
	}
}

func (r *TeamInvitationResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *TeamInvitationResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan TeamInvitationModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	invitation, err := r.client.CreateInvitation(ctx, plan.TenantID.ValueString(), client.CreateInvitationRequest{
		Email: plan.Email.ValueString(),
		Role:  plan.Role.ValueString(),
	})
	if err != nil {
		resp.Diagnostics.AddError("Error creating team invitation", err.Error())
		return
	}

	plan.ID = types.StringValue(invitation.ID)
	plan.TenantID = types.StringValue(invitation.TenantID)
	plan.Email = types.StringValue(invitation.Email)
	plan.Role = types.StringValue(invitation.Role)
	plan.Token = types.StringValue(invitation.Token)
	plan.InvitedBy = types.StringValue(invitation.InvitedBy)
	plan.ExpiresAt = types.StringValue(invitation.ExpiresAt)
	plan.CreatedAt = types.StringValue(invitation.CreatedAt)
	if invitation.AcceptedAt != nil {
		plan.AcceptedAt = types.StringValue(*invitation.AcceptedAt)
	} else {
		plan.AcceptedAt = types.StringNull()
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *TeamInvitationResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state TeamInvitationModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	invitations, err := r.client.ListInvitations(ctx, state.TenantID.ValueString())
	if err != nil {
		if client.IsNotFound(err) {
			tflog.Warn(ctx, "Team not found, removing invitation from state", map[string]interface{}{
				"tenant_id": state.TenantID.ValueString(),
			})
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError("Error reading team invitation", err.Error())
		return
	}

	var found bool
	for _, inv := range invitations {
		if inv.ID == state.ID.ValueString() {
			state.Email = types.StringValue(inv.Email)
			state.Role = types.StringValue(inv.Role)
			state.Token = types.StringValue(inv.Token)
			state.InvitedBy = types.StringValue(inv.InvitedBy)
			state.ExpiresAt = types.StringValue(inv.ExpiresAt)
			state.CreatedAt = types.StringValue(inv.CreatedAt)
			if inv.AcceptedAt != nil {
				state.AcceptedAt = types.StringValue(*inv.AcceptedAt)
			} else {
				state.AcceptedAt = types.StringNull()
			}
			found = true
			break
		}
	}

	if !found {
		tflog.Warn(ctx, "Team invitation not found, removing from state", map[string]interface{}{
			"id": state.ID.ValueString(),
		})
		resp.State.RemoveResource(ctx)
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *TeamInvitationResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// All user-settable fields use RequiresReplace, so Update should never be called.
	resp.Diagnostics.AddError(
		"Update Not Supported",
		"Team invitations are immutable. All changes require replacement.",
	)
}

func (r *TeamInvitationResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state TeamInvitationModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Warn(ctx, "Team invitation deletion is not supported by the API. Removing invitation from Terraform state only.", map[string]interface{}{
		"id":        state.ID.ValueString(),
		"tenant_id": state.TenantID.ValueString(),
		"email":     state.Email.ValueString(),
	})
}

func (r *TeamInvitationResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	parts := strings.SplitN(req.ID, "/", 2)
	if len(parts) != 2 || parts[0] == "" || parts[1] == "" {
		resp.Diagnostics.AddError(
			"Invalid Import ID",
			"Import ID must be in the format 'tenant_id/invitation_id'.",
		)
		return
	}

	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("tenant_id"), parts[0])...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), parts[1])...)
}
