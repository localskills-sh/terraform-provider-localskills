package user_token

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/localskills/terraform-provider-localskills/internal/client"
)

var (
	_ resource.Resource                = &userTokenResource{}
	_ resource.ResourceWithImportState = &userTokenResource{}
)

type userTokenResource struct {
	client *client.Client
}

func NewResource() resource.Resource {
	return &userTokenResource{}
}

func (r *userTokenResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_user_token"
}

func (r *userTokenResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Manages a user API token on localskills.sh.\n\n~> **Important:** The `token_value` attribute is only available at creation time. After creation, the API only stores a hash. If you lose the Terraform state, the token **cannot** be recovered.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: "The unique identifier of the token.",
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"name": schema.StringAttribute{
				Description: "The name of the token.",
				Required:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
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

func (r *userTokenResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *userTokenResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan UserTokenModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	token, err := r.client.CreateUserToken(ctx, client.CreateTokenRequest{
		Name: plan.Name.ValueString(),
	})
	if err != nil {
		resp.Diagnostics.AddError("Error creating user token", err.Error())
		return
	}

	plan.ID = types.StringValue(token.ID)
	plan.Name = types.StringValue(token.Name)
	plan.TokenValue = types.StringValue(token.Token)
	plan.LastUsedAt = types.StringNull()
	plan.ExpiresAt = types.StringNull()
	plan.CreatedAt = types.StringNull()

	// Read back to populate all fields
	tokens, err := r.client.ListUserTokens(ctx)
	if err == nil {
		for _, t := range tokens {
			if t.ID == token.ID {
				if t.LastUsedAt != nil {
					plan.LastUsedAt = types.StringValue(*t.LastUsedAt)
				}
				if t.ExpiresAt != nil {
					plan.ExpiresAt = types.StringValue(*t.ExpiresAt)
				}
				plan.CreatedAt = types.StringValue(t.CreatedAt)
				break
			}
		}
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *userTokenResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var currentState UserTokenModel
	resp.Diagnostics.Append(req.State.Get(ctx, &currentState)...)
	if resp.Diagnostics.HasError() {
		return
	}

	tokens, err := r.client.ListUserTokens(ctx)
	if err != nil {
		resp.Diagnostics.AddError("Error reading user tokens", err.Error())
		return
	}

	var found *client.ApiToken
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

	var state UserTokenModel
	state.ID = types.StringValue(found.ID)
	state.Name = types.StringValue(found.Name)
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

func (r *userTokenResource) Update(_ context.Context, _ resource.UpdateRequest, resp *resource.UpdateResponse) {
	resp.Diagnostics.AddError(
		"Update Not Supported",
		"User tokens cannot be updated. All attributes require replacement.",
	)
}

func (r *userTokenResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state UserTokenModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	err := r.client.DeleteUserToken(ctx, state.ID.ValueString())
	if err != nil {
		if client.IsNotFound(err) {
			return
		}
		resp.Diagnostics.AddError("Error deleting user token", err.Error())
	}
}

func (r *userTokenResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
