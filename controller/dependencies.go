package controller

import (
	"errors"

	"app/apierrors"
	"app/components/sms"
	"app/models"
	"app/pkg/validator"
	"app/utilities"

	"github.com/labstack/echo/v4"
	js "github.com/santhosh-tekuri/jsonschema/v5"
)

// RouterDependencies centralized router grouping dependencies
type RouterDependencies struct {
	// E echo router group usually the /api/v1 group
	E *echo.Group
	// Requires is a middleware that strictly requires a scope to be present
	Requires func(scope ...string) echo.MiddlewareFunc
}

// Dependencies centralized controller dependencies values are injected into
// each controller for easier maintenance
type Dependencies struct {
	// Schemas contains a map of each model/attribute json-schema
	// the key is represented as a url:
	//  config.DOMAIN + "/public/schemas/model-name.json"
	//  config.DOMAIN + "/public/schemas/properties/property.json"
	Schemas map[string]*js.Schema

	// Utils mostly handles the request context information
	Utils *utilities.Utils

	// APIErr is the centralized error handler and logger
	APIErr *apierrors.APIErrors

	// Models application data access layer
	// contains db connection and query builder instance
	// each model represents a database table
	Models *models.Models

	// SMS is used to notify customers
	SMS *sms.SMS
}

// GetValidator starts a validator instance with the selected schema name
// it returns error as APIErr.BadRequest format
func (d *Dependencies) GetValidator(
	ctx echo.Context,
	schemaName string,
) (*validator.Validator, error) {
	v, err := validator.NewValidator(
		ctx,
		d.Models.DB,
		d.Utils.Logger,
		d.Schemas,
		schemaName,
	)
	if err != nil {
		return nil, d.APIErr.BadRequest(ctx, err)
	}
	if v == nil {
		err := errors.New("null validator")
		return nil, d.APIErr.BadRequest(ctx, err)
	}
	return v, nil
}
