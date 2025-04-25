package resources

import (
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
)

func AllowlistsSchema() schema.Schema {
	return schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id": stringAttributeSchema(
				WithComputed(),
				WithUseStateForUnknown(),
			),
			"organization_id": stringAttributeSchema(
				WithRequired(),
				WithRequiresReplace(),
				WithDescription("The organization id the allowlist belongs to"),
			),
			"project_id": stringAttributeSchema(
				WithRequired(),
				WithRequiresReplace(),
				WithDescription("The Project id the allowlist belongs to"),
			),
			"cluster_id": stringAttributeSchema(
				WithRequired(),
				WithRequiresReplace(),
				WithDescription("The cluster id the allowlist belongs to"),
			),
			"cidr": stringAttributeSchema(
				WithRequired(),
				WithRequiresReplace(),
				WithDescription("The IP CIDR range to allow access"),
			),
			"comment": stringAttributeSchema(
				WithOptional(),
				WithRequiresReplace(),
				WithDescription("A comment/description about the allow list"),
				WithComputed(),
			),
			"expires_at": stringAttributeSchema(
				WithOptional(),
				WithRequiresReplace(),
				WithDescription("The expiration time of the allow list"),
			),
			"audit": computedAuditAttribute(),
		},
	}
}
