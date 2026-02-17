package team_invitations

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/localskills/terraform-provider-localskills/internal/client"
)

var (
	_ datasource.DataSource              = &TeamInvitationsDataSource{}
	_ datasource.DataSourceWithConfigure = &TeamInvitationsDataSource{}
)

type TeamInvitationsDataSource struct {
	client *client.Client
}

func NewDataSource() datasource.DataSource {
	return &TeamInvitationsDataSource{}
}

func (d *TeamInvitationsDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_team_invitations"
}

func (d *TeamInvitationsDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Fetches all invitations for a team.",
		Attributes: map[string]schema.Attribute{
			"tenant_id": schema.StringAttribute{
				Description: "The ID of the team (tenant) to list invitations for.",
				Required:    true,
			},
			"invitations": schema.ListNestedAttribute{
				Description: "List of invitations.",
				Computed:    true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.StringAttribute{
							Description: "The unique identifier of the invitation.",
							Computed:    true,
						},
						"tenant_id": schema.StringAttribute{
							Description: "The ID of the team.",
							Computed:    true,
						},
						"email": schema.StringAttribute{
							Description: "The email address of the invited person.",
							Computed:    true,
						},
						"role": schema.StringAttribute{
							Description: "The role assigned to the invited user.",
							Computed:    true,
						},
						"invited_by": schema.StringAttribute{
							Description: "The ID of the user who created the invitation.",
							Computed:    true,
						},
						"expires_at": schema.StringAttribute{
							Description: "The timestamp when the invitation expires.",
							Computed:    true,
						},
						"accepted_at": schema.StringAttribute{
							Description: "The timestamp when the invitation was accepted.",
							Computed:    true,
						},
						"created_at": schema.StringAttribute{
							Description: "The timestamp when the invitation was created.",
							Computed:    true,
						},
					},
				},
			},
		},
	}
}

func (d *TeamInvitationsDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *TeamInvitationsDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var config TeamInvitationsDataSourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)
	if resp.Diagnostics.HasError() {
		return
	}

	invitations, err := d.client.ListInvitations(ctx, config.TenantID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Error listing team invitations", err.Error())
		return
	}

	var state TeamInvitationsDataSourceModel
	state.TenantID = config.TenantID

	for _, inv := range invitations {
		model := InvitationModel{
			ID:        types.StringValue(inv.ID),
			TenantID:  types.StringValue(inv.TenantID),
			Email:     types.StringValue(inv.Email),
			Role:      types.StringValue(inv.Role),
			InvitedBy: types.StringValue(inv.InvitedBy),
			ExpiresAt: types.StringValue(inv.ExpiresAt),
			CreatedAt: types.StringValue(inv.CreatedAt),
		}
		if inv.AcceptedAt != nil {
			model.AcceptedAt = types.StringValue(*inv.AcceptedAt)
		} else {
			model.AcceptedAt = types.StringNull()
		}
		state.Invitations = append(state.Invitations, model)
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}
