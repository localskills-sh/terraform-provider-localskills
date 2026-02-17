package user_profile

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/localskills/terraform-provider-localskills/internal/client"
)

var _ datasource.DataSource = &userProfileDataSource{}

type userProfileDataSource struct {
	client *client.Client
}

func NewDataSource() datasource.DataSource {
	return &userProfileDataSource{}
}

func (d *userProfileDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_user_profile"
}

func (d *userProfileDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Fetches the current user's profile.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: "The unique identifier of the user.",
				Computed:    true,
			},
			"username": schema.StringAttribute{
				Description: "The username of the user.",
				Computed:    true,
			},
			"name": schema.StringAttribute{
				Description: "The display name of the user.",
				Computed:    true,
			},
			"email": schema.StringAttribute{
				Description: "The email address of the user.",
				Computed:    true,
			},
			"image": schema.StringAttribute{
				Description: "The avatar image URL of the user.",
				Computed:    true,
			},
			"bio": schema.StringAttribute{
				Description: "The bio of the user.",
				Computed:    true,
			},
		},
	}
}

func (d *userProfileDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *userProfileDataSource) Read(ctx context.Context, _ datasource.ReadRequest, resp *datasource.ReadResponse) {
	profile, err := d.client.GetUserProfile(ctx)
	if err != nil {
		resp.Diagnostics.AddError("Error reading user profile", err.Error())
		return
	}

	var state UserProfileModel
	state.ID = types.StringValue(profile.ID)
	state.Email = types.StringValue(profile.Email)

	if profile.Username != nil {
		state.Username = types.StringValue(*profile.Username)
	} else {
		state.Username = types.StringNull()
	}
	if profile.Name != nil {
		state.Name = types.StringValue(*profile.Name)
	} else {
		state.Name = types.StringNull()
	}
	if profile.Image != nil {
		state.Image = types.StringValue(*profile.Image)
	} else {
		state.Image = types.StringNull()
	}
	if profile.Bio != nil {
		state.Bio = types.StringValue(*profile.Bio)
	} else {
		state.Bio = types.StringNull()
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}
