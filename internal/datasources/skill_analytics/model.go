package skill_analytics

import (
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type SkillAnalyticsDataSourceModel struct {
	SkillID        types.String `tfsdk:"skill_id"`
	TotalDownloads types.Int64  `tfsdk:"total_downloads"`
	UniqueUsers    types.Int64  `tfsdk:"unique_users"`
	UniqueIPs      types.Int64  `tfsdk:"unique_ips"`
}
