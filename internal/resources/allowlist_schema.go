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
			"organization_id": stringAttributeWithDescription([]string{required, requiresReplace}, "The organization ID"),
			"project_id":      stringAttributeWithDescription([]string{required, requiresReplace}, "The project ID"),
			"cluster_id":      stringAttributeWithDescription([]string{required, requiresReplace}, "The cluster ID"),
			"cidr":            stringAttributeWithDescription([]string{required, requiresReplace}, "A cidr range for allowed ip's"),
			"comment":         stringAttributeWithDescription([]string{optional, requiresReplace, computed}, "A comment/description about the allow list"),
			"expires_at":      stringAttributeWithDescription([]string{optional, requiresReplace}, "An expires at timestamp as a string"),
			"audit":           computedAuditAttribute(),
		},
	}
}
