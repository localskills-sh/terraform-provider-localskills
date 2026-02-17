package team_token

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/localskills/terraform-provider-localskills/internal/client"
)

var (
	_ resource.Resource                = &teamTokenResource{}
	_ resource.ResourceWithImportState = &teamTokenResource{}
)

type teamTokenResource struct {
	client *client.Client
}

func NewResource() resource.Resource {
	return &teamTokenResource{}
}

func (r *teamTokenResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_team_token"
}

func (r *teamTokenResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Manages a team API token on localskills.sh.\n\n~> **Important:** The `token_value` attribute is only available at creation time. After creation, the API only stores a hash. If you lose the Terraform state, the token **cannot** be recovered.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: "The unique identifier of the token.",
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"tenant_id": schema.StringAttribute{
				Description: "The tenant (team) ID.",
				Required:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"name": schema.StringAttribute{
				Description: "The name of the token.",
				Required:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"expires_in_days": schema.Int64Attribute{
				Description: "Number of days until the token expires.",
				Optional:    true,
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.RequiresReplace(),
				},
			},
			"token_value": schema.StringAttribute{
				Description: "The secret token value. Only available after creation.",
				Computed:    true,
				Sensitive:   true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"last_used_at": schema.StringAttribute{
				Description: "When the token was last used.",
				Computed:    true,
			},
			"expires_at": schema.StringAttribute{
				Description: "When the token expires.",
				Computed:    true,
			},
			"created_at": schema.StringAttribute{
				Description: "When the token was created.",
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
		},
	}
}

func (r *teamTokenResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *teamTokenResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan TeamTokenModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	createReq := client.CreateTeamTokenRequest{
		Name: plan.Name.ValueString(),
	}
	if !plan.ExpiresInDays.IsNull() && !plan.ExpiresInDays.IsUnknown() {
		days := int(plan.ExpiresInDays.ValueInt64())
		createReq.ExpiresInDays = &days
	}

	token, err := r.client.CreateTeamToken(ctx, plan.TenantID.ValueString(), createReq)
	if err != nil {
		resp.Diagnostics.AddError("Error creating team token", err.Error())
		return
	}

	plan.ID = types.StringValue(token.ID)
	plan.Name = types.StringValue(token.Name)
	plan.TokenValue = types.StringValue(token.Token)
	plan.CreatedAt = types.StringValue(token.CreatedAt)
	if token.ExpiresAt != nil {
		plan.ExpiresAt = types.StringValue(*token.ExpiresAt)
	} else {
		plan.ExpiresAt = types.StringNull()
	}
	plan.LastUsedAt = types.StringNull()

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *teamTokenResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var currentState TeamTokenModel
	resp.Diagnostics.Append(req.State.Get(ctx, &currentState)...)
	if resp.Diagnostics.HasError() {
		return
	}

	tokens, err := r.client.ListTeamTokens(ctx, currentState.TenantID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Error reading team tokens", err.Error())
		return
	}

	var found *client.TeamApiToken
	for _, t := range tokens {
		if t.ID == currentState.ID.ValueString() {
			found = &t
			break
		}
	}

	if found == nil {
		resp.State.RemoveResource(ctx)
		return
	}

	var state TeamTokenModel
	state.ID = types.StringValue(found.ID)
	state.TenantID = currentState.TenantID
	state.Name = types.StringValue(found.Name)
	state.ExpiresInDays = currentState.ExpiresInDays
	if found.LastUsedAt != nil {
		state.LastUsedAt = types.StringValue(*found.LastUsedAt)
	} else {
		state.LastUsedAt = types.StringNull()
	}
	if found.ExpiresAt != nil {
		state.ExpiresAt = types.StringValue(*found.ExpiresAt)
	} else {
		state.ExpiresAt = types.StringNull()
	}
	state.CreatedAt = types.StringValue(found.CreatedAt)

	// Preserve token_value from state since API only returns hashes
	state.TokenValue = currentState.TokenValue

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *teamTokenResource) Update(_ context.Context, _ resource.UpdateRequest, resp *resource.UpdateResponse) {
	resp.Diagnostics.AddError(
		"Update Not Supported",
		"Team tokens cannot be updated. All attributes require replacement.",
	)
}

func (r *teamTokenResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state TeamTokenModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	err := r.client.DeleteTeamToken(ctx, state.TenantID.ValueString(), state.ID.ValueString())
	if err != nil {
		if client.IsNotFound(err) {
			return
		}
		resp.Diagnostics.AddError("Error deleting team token", err.Error())
	}
}

func (r *teamTokenResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	parts := strings.SplitN(req.ID, "/", 2)
	if len(parts) != 2 || parts[0] == "" || parts[1] == "" {
		resp.Diagnostics.AddError(
			"Invalid Import ID",
			"Expected import ID in the format: tenant_id/token_id",
		)
		return
	}

	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("tenant_id"), parts[0])...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), parts[1])...)
}
