package resources

import (
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
)

func AllowlistsSchema() schema.Schema {
	return schema.Schema{
		Description: "Couchbase Capella only allows trusted IP addresses to connect to databases. Each database has a configurable Allowed IP list that can include up to 75 entries. Each entry can be a single IP address or an IP address space. Any IP address you add to this list can have a user-specified expiration time for temporary access, or be permanent. Capella automatically denies any connection attempts to and from an IP not in the allowed IP list.",
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
