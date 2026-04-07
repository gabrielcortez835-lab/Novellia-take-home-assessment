package fhir

import (
	"fmt"
	"strings"

	"github.com/gabrielcortez835-lab/Novellia-take-home-assessment/constants"
	"github.com/tidwall/gjson"
)

func ValidateResource(resourceType string, json gjson.Result) []string {
	var warnings []string

	switch resourceType {
	case "Condition":
		warnings = validateConditionGJSON(json)
	case "Observation":
		warnings = validateObservationGJSON(json)
	case "MedicationRequest":
		warnings = validateMedicationRequestGJSON(json)
	case "Procedure":
		warnings = validateProcedureGJSON(json)
	case "Patient":
		warnings = validatePatientGJSON(json)
	default:
		warnings = []string{fmt.Sprintf("Unknown ResourceType: %s", resourceType)}
	}

	return warnings
}

func validateConditionGJSON(resource gjson.Result) []string {
	retArr := []string{}

	// FHIRBase check
	retArr = append(retArr, validateFHIRBaseGJSON(resource)...)

	// Subject check
	retArr = append(retArr, validateSubjectGJSON(resource.Get("subject"))...)

	return retArr
}

func validateObservationGJSON(resource gjson.Result) []string {
	retArr := []string{}

	retArr = append(retArr, validateFHIRBaseGJSON(resource)...)

	if err := validateStatusGJSON(resource.Get("status")); err != "" {
		retArr = append(retArr, err)
	}

	retArr = append(retArr, validateCodeableConceptGJSON(resource.Get("code"))...)
	retArr = append(retArr, validateSubjectGJSON(resource.Get("subject"))...)

	return retArr
}

func validateMedicationRequestGJSON(resource gjson.Result) []string {
	retArr := []string{}

	retArr = append(retArr, validateFHIRBaseGJSON(resource)...)

	if err := validateStatusGJSON(resource.Get("status")); err != "" {
		retArr = append(retArr, err)
	}

	retArr = append(retArr, validateCodeableConceptGJSON(resource.Get("medicationCodeableConcept"))...)
	retArr = append(retArr, validateSubjectGJSON(resource.Get("subject"))...)

	return retArr
}

func validateProcedureGJSON(resource gjson.Result) []string {
	retArr := []string{}

	retArr = append(retArr, validateFHIRBaseGJSON(resource)...)
	retArr = append(retArr, validateSubjectGJSON(resource.Get("subject"))...)

	return retArr
}

func validatePatientGJSON(resource gjson.Result) []string {
	retArr := []string{}
	retArr = append(retArr, validateFHIRBaseGJSON(resource)...)

	return retArr
}

func validateFHIRBaseGJSON(resource gjson.Result) []string {
	retArr := []string{}

	id := strings.TrimSpace(resource.Get("id").String())
	resourceType := strings.TrimSpace(resource.Get("resourceType").String())

	if id == "" {
		retArr = append(retArr, "Missing field: id")
	}
	if !constants.ValidResourceTypes[resourceType] {
		retArr = append(retArr, fmt.Sprintf("Invalid ResourceType: %s", resourceType))
	}

	return retArr
}

func validateStatusGJSON(status gjson.Result) string {
	statusStr := strings.TrimSpace(status.String())
	if !constants.ValidStatus[statusStr] {
		return fmt.Sprintf("Invalid Status: %s", statusStr)
	}
	return ""
}

func validateCodeableConceptGJSON(codeableConcept gjson.Result) []string {
	retArr := []string{}

	codingArr := codeableConcept.Get("coding").Array()
	if len(codingArr) == 0 {
		retArr = append(retArr, "Missing Field: Coding")
	}
	for i, coding := range codingArr {
		retArr = append(retArr, validateCodingGJSON(i, coding)...)
	}

	return retArr
}

func validateCodingGJSON(pos int, coding gjson.Result) []string {
	retArr := []string{}

	if strings.TrimSpace(coding.Get("system").String()) == "" {
		retArr = append(retArr, fmt.Sprintf("Coding[%d] missing field: System", pos))
	}
	if strings.TrimSpace(coding.Get("code").String()) == "" {
		retArr = append(retArr, fmt.Sprintf("Coding[%d] missing field: Code", pos))
	}
	if strings.TrimSpace(coding.Get("display").String()) == "" {
		retArr = append(retArr, fmt.Sprintf("Coding[%d] missing field: Display", pos))
	}

	return retArr
}

func validateSubjectGJSON(subject gjson.Result) []string {
	retArr := []string{}

	if strings.TrimSpace(subject.Get("reference").String()) == "" {
		retArr = append(retArr, "Subject missing field: Reference")
	}
	if strings.TrimSpace(subject.Get("display").String()) == "" {
		retArr = append(retArr, "Subject missing field: Display")
	}

	return retArr
}
