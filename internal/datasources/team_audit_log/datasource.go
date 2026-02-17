package team_audit_log

import (
	"context"
	"fmt"
	"strconv"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/localskills/terraform-provider-localskills/internal/client"
)

var _ datasource.DataSource = &teamAuditLogDataSource{}

type teamAuditLogDataSource struct {
	client *client.Client
}

func NewDataSource() datasource.DataSource {
	return &teamAuditLogDataSource{}
}

func (d *teamAuditLogDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_team_audit_log"
}

func (d *teamAuditLogDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Fetches the audit log for a team on localskills.sh. Supports pagination with `page` and `limit`, and filtering by `action` type.",
		Attributes: map[string]schema.Attribute{
			"tenant_id": schema.StringAttribute{
				Description: "The ID of the team (tenant) to fetch audit logs for.",
				Required:    true,
			},
			"page": schema.Int64Attribute{
				Description: "Page number to fetch.",
				Optional:    true,
			},
			"limit": schema.Int64Attribute{
				Description: "Number of entries per page.",
				Optional:    true,
			},
			"action": schema.StringAttribute{
				Description: "Filter by action type.",
				Optional:    true,
			},
			"total": schema.Int64Attribute{
				Description: "Total number of audit log entries.",
				Computed:    true,
			},
			"entries": schema.ListNestedAttribute{
				Description: "The audit log entries.",
				Computed:    true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.StringAttribute{
							Description: "The unique identifier of the audit log entry.",
							Computed:    true,
						},
						"action": schema.StringAttribute{
							Description: "The action that was performed.",
							Computed:    true,
						},
						"actor_id": schema.StringAttribute{
							Description: "The ID of the actor who performed the action.",
							Computed:    true,
						},
						"actor_name": schema.StringAttribute{
							Description: "The name of the actor.",
							Computed:    true,
						},
						"actor_image": schema.StringAttribute{
							Description: "The avatar image URL of the actor.",
							Computed:    true,
						},
						"resource_type": schema.StringAttribute{
							Description: "The type of resource affected.",
							Computed:    true,
						},
						"resource_id": schema.StringAttribute{
							Description: "The ID of the resource affected.",
							Computed:    true,
						},
						"metadata": schema.StringAttribute{
							Description: "Additional metadata as JSON.",
							Computed:    true,
						},
						"created_at": schema.StringAttribute{
							Description: "When the action was performed.",
							Computed:    true,
						},
					},
				},
			},
		},
	}
}

func (d *teamAuditLogDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *teamAuditLogDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data TeamAuditLogModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	params := map[string]string{}
	if !data.Page.IsNull() && !data.Page.IsUnknown() {
		params["page"] = strconv.FormatInt(data.Page.ValueInt64(), 10)
	}
	if !data.Limit.IsNull() && !data.Limit.IsUnknown() {
		params["limit"] = strconv.FormatInt(data.Limit.ValueInt64(), 10)
	}
	if !data.Action.IsNull() && !data.Action.IsUnknown() {
		params["action"] = data.Action.ValueString()
	}

	result, err := d.client.ListTeamAuditLog(ctx, data.TenantID.ValueString(), params)
	if err != nil {
		resp.Diagnostics.AddError("Error reading team audit log", err.Error())
		return
	}

	data.Total = types.Int64Value(int64(result.Total))
	data.Entries = make([]AuditLogEntryModel, len(result.Entries))
	for i, e := range result.Entries {
		entry := AuditLogEntryModel{
			ID:           types.StringValue(e.ID),
			Action:       types.StringValue(e.Action),
			ResourceType: types.StringValue(e.ResourceType),
			ResourceID:   types.StringValue(e.ResourceID),
			Metadata:     types.StringValue(e.Metadata),
			CreatedAt:    types.StringValue(e.CreatedAt),
		}
		if e.ActorID != nil {
			entry.ActorID = types.StringValue(*e.ActorID)
		} else {
			entry.ActorID = types.StringNull()
		}
		if e.ActorName != nil {
			entry.ActorName = types.StringValue(*e.ActorName)
		} else {
			entry.ActorName = types.StringNull()
		}
		if e.ActorImage != nil {
			entry.ActorImage = types.StringValue(*e.ActorImage)
		} else {
			entry.ActorImage = types.StringNull()
		}
		data.Entries[i] = entry
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
