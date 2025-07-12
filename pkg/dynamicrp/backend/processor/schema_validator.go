package processor

import (
	"context"
	"errors"
	"fmt"

	"github.com/radius-project/radius/pkg/ucp/api/v20231001preview"
	"github.com/radius-project/radius/pkg/ucp/resources"
)

var ErrNoSchemaFound = errors.New("no schema found for resource type")

/*
// SchemaValidator handles resource schema validation
type SchemaValidator struct {
	schemaValidator *schema.Validator
}

// NewSchemaValidator creates a new schema validator instance
func NewSchemaValidator() *SchemaValidator {
	return &SchemaValidator{
		schemaValidator: schema.NewValidator(),
	}
}
*/

// GetSchemaForResourceType fetches the schema for a resource type from UCP
func GetSchemaForResourceType(ctx context.Context, ucp *v20231001preview.ClientFactory, resourceType string, apiVersion string) (interface{}, error) {
	// Parse resource type to get components
	parsed, err := resources.ParseResource(resourceType)
	if err != nil {
		return nil, fmt.Errorf("invalid resource type format: %w", err)
	}

	// Fetch the API version resource to get schema
	providerNamespace := parsed.ProviderNamespace()
	resourceTypeName := parsed.Type()
	planeName := "radius" // Default plane name

	response, err := ucp.NewAPIVersionsClient().Get(
		ctx,
		planeName,
		providerNamespace,
		resourceTypeName,
		apiVersion,
		nil)
	if err != nil {
		return nil, ErrNoSchemaFound
	}

	// Check if schema exists
	if response.Properties == nil || response.Properties.Schema == nil {
		return nil, ErrNoSchemaFound
	}

	return response.Properties.Schema, nil
}

/*
// ResourceValidationError represents validation errors for a resource
type ResourceValidationError struct {
	ResourceType string
	APIVersion   string
	Errors       []schema.ValidationError
}

// Error returns the error message
func (e *ResourceValidationError) Error() string {
	return fmt.Sprintf("validation failed for resource type %s (apiVersion: %s): %d errors found",
		e.ResourceType, e.APIVersion, len(e.Errors))
}


// DetailedErrors returns detailed validation errors
func (e *ResourceValidationError) DetailedErrors() []schema.ValidationError {
	return e.Errors
}


// ValidateResource validates a resource JSON against a schema JSON
func (v *SchemaValidator) ValidateResource(resourceJSON []byte, schemaJSON []byte) error {
	// Parse resource JSON
	var resourceData map[string]interface{}
	if err := json.Unmarshal(resourceJSON, &resourceData); err != nil {
		return fmt.Errorf("failed to parse resource JSON: %w", err)
	}

	// Parse schema JSON to OpenAPI schema
	var schemaDataRaw interface{}
	if err := json.Unmarshal(schemaJSON, &schemaDataRaw); err != nil {
		return fmt.Errorf("failed to parse schema JSON: %w", err)
	}

	// Validate the resource data against the schema using the schema package
	if err := schema.ValidateResourceAgainstSchema(resourceData, schemaDataRaw); err != nil {
		return fmt.Errorf("validation failed: %w", err)
	}

	return nil
}
*/
