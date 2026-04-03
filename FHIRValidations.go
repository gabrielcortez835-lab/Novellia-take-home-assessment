package main

import (
	"fmt"
	"strings"
)

func validateCondition(resource Condition) []string {
	retArr := make([]string, 0)
	validationErrors := validateFHIRBase(resource.ID, resource.ResourceType)
	if len(validationErrors) > 0 {
		retArr = append(retArr, validationErrors...)
	}

	if validationErrors := validateSubject(resource.Subject); len(validationErrors) > 0 {
		retArr = append(retArr, validationErrors...)
	}

	return retArr
}

func validateObservation(resource Observation) []string {
	retArr := make([]string, 0)

	if validationErrors := validateFHIRBase(resource.ID, resource.ResourceType); len(validationErrors) > 0 {
		retArr = append(retArr, validationErrors...)
	}

	if validationError := validateStatus(resource.Status); validationError != "" {
		retArr = append(retArr, validationError)
	}

	if validationErrors := validateCodeableConcept(resource.Code); len(validationErrors) > 0 {
		retArr = append(retArr, validationErrors...)
	}

	if validationErrors := validateSubject(resource.Subject); len(validationErrors) > 0 {
		retArr = append(retArr, validationErrors...)
	}

	return retArr
}

func validateMedicationRequest(resource MedicationRequest) []string {
	retArr := make([]string, 0)

	if validationErrors := validateFHIRBase(resource.ID, resource.ResourceType); len(validationErrors) > 0 {
		retArr = append(retArr, validationErrors...)
	}

	if validationError := validateStatus(resource.Status); validationError != "" {
		retArr = append(retArr, validationError)
	}

	if validationErrors := validateCodeableConcept(resource.MedicationCodeableConcept); len(validationErrors) > 0 {
		retArr = append(retArr, validationErrors...)
	}

	if validationErrors := validateSubject(resource.Subject); len(validationErrors) > 0 {
		retArr = append(retArr, validationErrors...)
	}

	return retArr
}

func validateProcedure(resource Procedure) []string {
	retArr := make([]string, 0)
	validationErrors := validateFHIRBase(resource.ID, resource.ResourceType)
	if len(validationErrors) > 0 {
		retArr = append(retArr, validationErrors...)
	}

	if validationErrors := validateSubject(resource.Subject); len(validationErrors) > 0 {
		retArr = append(retArr, validationErrors...)
	}

	return retArr
}

func validatePatient(resource Patient) []string {
	retArr := make([]string, 0)
	validationErrors := validateFHIRBase(resource.ID, resource.ResourceType)
	if len(validationErrors) > 0 {
		retArr = append(retArr, validationErrors...)
	}

	return retArr
}

func validateFHIRBase(id string, resourceType string) []string {
	retArr := make([]string, 0)

	if strings.TrimSpace(id) == "" {
		retArr = append(retArr, "missing field: id")
	}

	if !ValidResourceTypes[resourceType] {
		retArr = append(retArr, fmt.Sprintf("Invalid ResourceType: %s", resourceType))
	}

	return retArr
}

func validateStatus(status string) string {
	if !ValidStatus[status] {
		return fmt.Sprintf("Invalid Status: %s", status)
	}
	return ""
}

func validateCodeableConcept(codeableConcept CodeableConcept) []string {
	retArr := make([]string, 0)
	if len(codeableConcept.Coding) == 0 {
		retArr = append(retArr, "Missing Field: Coding")
	}

	for i := range codeableConcept.Coding {
		if codingErrors := validateCoding(i, codeableConcept.Coding[i]); len(codingErrors) > 0 {
			retArr = append(retArr, codingErrors...)
		}
	}
	return retArr
}

func validateCoding(pos int, coding Coding) []string {
	retArr := make([]string, 0)
	if strings.TrimSpace(coding.System) == "" {
		retArr = append(retArr, fmt.Sprintf("Coding[%d] missing field: System", pos))
	}
	if strings.TrimSpace(coding.Code) == "" {
		retArr = append(retArr, fmt.Sprintf("Coding[%d] missing field: Code", pos))
	}
	if strings.TrimSpace(coding.Display) == "" {
		retArr = append(retArr, fmt.Sprintf("Coding[%d] missing field: Display", pos))
	}

	return retArr
}

func validateSubject(subject Reference) []string {
	retArr := make([]string, 0)
	if strings.TrimSpace(subject.Reference) == "" {
		retArr = append(retArr, "Subject missing field: Reference")
	}
	if strings.TrimSpace(subject.Display) == "" {
		retArr = append(retArr, "Subject missing field: Display")
	}

	return retArr
}
