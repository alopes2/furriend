// Pacakge domain containing domain libraries
package domain

import (
	"slices"
	"strings"
)

// ValidationResult is result of a validation
// IsValid is true when validation succeeds
// if false, Errors will have values
type ValidationResult struct {
	IsValid bool
	Errors  []string
}

const (
	petNameMaxLength = 50
)

// PetRequestValidator validates a Pet from the request
func PetRequestValidator(pet Pet) ValidationResult {
	validationResult := ValidationResult{}
	petSpecies := []string{
		"DOG",
		"CAT",
	}

	if strings.TrimSpace(pet.Name) == "" {
		validationResult.Errors = append(validationResult.Errors, "Pet name required")
	}

	if len(pet.Name) > petNameMaxLength {
		validationResult.Errors = append(validationResult.Errors, "Pet name should be less than 50 chars")
	}

	if strings.TrimSpace(pet.Specie) == "" {
		validationResult.Errors = append(validationResult.Errors, "Pet specie required")
	}

	if !slices.Contains(petSpecies, strings.ToUpper(pet.Specie)) {
		validationResult.Errors = append(validationResult.Errors, "Pet specie invalid")
	}

	validationResult.IsValid = len(validationResult.Errors) == 0

	return validationResult
}
