package validation

import (
	"context"

	"github.com/go-logr/logr"
)

type ResourceValidator interface {
	Validate(ctx context.Context) error
	GetResource(ctx context.Context) (interface{}, error)
	Name() string
}

type ValidationManager struct {
	logger         *logr.Logger
	validators     []ResourceValidator
	validatorNames map[string]bool
}
