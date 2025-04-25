package resources

import (
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
)

type StringAttributeOption func(*schema.StringAttribute)

func stringAttributeSchema(options ...StringAttributeOption) *schema.StringAttribute {
	attribute := &schema.StringAttribute{}

	for _, option := range options {
		option(attribute)
	}

	return attribute
}

func WithRequired() StringAttributeOption {
	return func(attr *schema.StringAttribute) {
		attr.Required = true
	}
}

func WithOptional() StringAttributeOption {
	return func(attr *schema.StringAttribute) {
		attr.Optional = true
	}
}

func WithComputed() StringAttributeOption {
	return func(attr *schema.StringAttribute) {
		attr.Computed = true
	}
}

func WithSensitive() StringAttributeOption {
	return func(attr *schema.StringAttribute) {
		attr.Sensitive = true
	}
}

func WithRequiresReplace() StringAttributeOption {
	return func(attr *schema.StringAttribute) {
		attr.PlanModifiers = append(attr.PlanModifiers, stringplanmodifier.RequiresReplace())
	}
}

func WithUseStateForUnknown() StringAttributeOption {
	return func(attr *schema.StringAttribute) {
		attr.PlanModifiers = append(attr.PlanModifiers, stringplanmodifier.UseStateForUnknown())
	}
}

func WithDeprecated(deprecationMessage string) StringAttributeOption {
	return func(attr *schema.StringAttribute) {
		attr.DeprecationMessage = deprecationMessage
	}
}

func WithValidators(validators ...validator.String) StringAttributeOption {
	return func(attr *schema.StringAttribute) {
		attr.Validators = append(attr.Validators, validators...)
	}
}

func WithDescription(description string) StringAttributeOption {
	return func(attr *schema.StringAttribute) {
		attr.Description = description
	}
}
