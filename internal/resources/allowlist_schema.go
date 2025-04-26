package resources

import (
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
)

func AllowlistsSchema() schema.Schema {
	return schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
				Description: "The ID of the allowlist.",
			},
			"organization_id": stringAttributeWithDescription([]string{required, requiresReplace}, "The GUID4 ID of the organization."),
			"project_id":      stringAttributeWithDescription([]string{required, requiresReplace}, "The GUID4 ID of the project."),
			"cluster_id":      stringAttributeWithDescription([]string{required, requiresReplace}, "The GUID4 ID of the cluster."),
			"cidr":            stringAttributeWithDescription([]string{required, requiresReplace}, "The trusted CIDR to allow the database connections from."),
			"comment":         stringAttributeWithDescription([]string{optional, requiresReplace, computed}, "A short description of the allowed CIDR."),
			"expires_at":      stringAttributeWithDescription([]string{optional, requiresReplace}, "An RFC3339 timestamp determining when the allowed CIDR should expire. If this field is empty/omitted then the allowed CIDR is permanent and will never automatically expire."),
			"audit":           computedAuditAttribute(),
		},
	}
}
