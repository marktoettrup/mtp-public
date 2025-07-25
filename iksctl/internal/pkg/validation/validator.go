// Package validation provides generic resource validation capabilities
package validation

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/go-logr/logr"
)

const (
	ExitCodeSuccess       = 0
	ExitCodeInternalError = 1
	ExitCodeBadInput      = 2
)

func NewValidationManager(logger *logr.Logger) *ValidationManager {
	return &ValidationManager{
		logger:         logger,
		validators:     []ResourceValidator{},
		validatorNames: make(map[string]bool),
	}
}

func (vm *ValidationManager) AddValidator(validator ResourceValidator) {
	name := validator.Name()
	if !vm.validatorNames[name] {
		vm.validators = append(vm.validators, validator)
		vm.validatorNames[name] = true
	}
}

func (vm *ValidationManager) ValidateAll(ctx context.Context) error {
	var errs []string

	resourceGroups := make(map[string][]ResourceValidator)
	for _, validator := range vm.validators {
		name := validator.Name()
		parts := strings.SplitN(name, "/", 2)
		if len(parts) == 2 {
			resourceType := parts[0]
			resourceGroups[resourceType] = append(resourceGroups[resourceType], validator)
		} else {
			resourceGroups["other"] = append(resourceGroups["other"], validator)
		}
	}

	for resourceType, validators := range resourceGroups {
		if len(validators) > 1 {
			vm.logger.Info(fmt.Sprintf("Validating %d %s resources", len(validators), resourceType))
		}

		for _, validator := range validators {
			vm.logger.Info("Validating resource", "resource", validator.Name())
			if err := validator.Validate(ctx); err != nil {
				errs = append(errs, fmt.Sprintf("%s: %s", validator.Name(), err.Error()))
			} else {
				vm.logger.Info("Resource validation successful", "resource", validator.Name())
			}
		}
	}

	if len(errs) > 0 {
		return errors.New(strings.Join(errs, ", "))
	}

	return nil
}
