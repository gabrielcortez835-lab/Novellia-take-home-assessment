package apiFunctions

import (
	"encoding/json"
	"fmt" // for printing and formatting

	"github.com/gabrielcortez835-lab/Novellia-take-home-assessment/constants"
	"github.com/gabrielcortez835-lab/Novellia-take-home-assessment/extractionConfig"
	"github.com/gabrielcortez835-lab/Novellia-take-home-assessment/fhir"
	"github.com/gabrielcortez835-lab/Novellia-take-home-assessment/sql"
	"github.com/tidwall/gjson"
)

func ProcessImportedJson(jsonString string, cfg extractionConfig.ExtractionConfig) ([]string, error) {
	// Parse once into a gjson.Result
	resource := gjson.Parse(jsonString)

	resourceType := resource.Get("resourceType").String()
	if resourceType == "" {
		return nil, fmt.Errorf("missing resourceType")
	}

	id := resource.Get("id").String()
	if id == "" {
		return nil, fmt.Errorf("missing id")
	}

	patientID := ""
	if resourceType == constants.ResourceTypeEnum.Patient {
		patientID = id
	} else {
		subj := resource.Get("subject.reference")
		if subj.Exists() {
			patientID = subj.String()
		}
	}

	extractedJson := extractFieldsToGJSON(resource, cfg)

	warnings := fhir.ValidateResource(resourceType, extractedJson)

	if len(warnings) > 0 {
		jsonBytes, err := json.Marshal(warnings)

		if err != nil {
			return nil, err
		}

		sql.InsertValidationError(id, string(jsonBytes))
	}

	err := sql.InsertResource(id, resourceType, patientID, jsonString)
	if err != nil {
		return nil, err
	}

	return warnings, nil
}

func extractFieldsToGJSON(resource gjson.Result, cfg extractionConfig.ExtractionConfig) gjson.Result {
	extracted := make(map[string]any)
	resourceType := resource.Get("resourceType").String()

	for fieldPath, allowed := range cfg.Fields {
		shouldExtract := false

		switch v := allowed.(type) {
		case string:
			if v == "all" {
				shouldExtract = true
			}

		case []interface{}:
			for _, t := range v {
				if s, ok := t.(string); ok && s == resourceType {
					shouldExtract = true
					break
				}
			}
		}

		if shouldExtract {
			val := resource.Get(fieldPath)
			if val.Exists() {
				extracted[fieldPath] = val.Value()
			}
		}
	}

	// Marshal the extracted object to a JSON string
	jsonBytes, err := json.Marshal(extracted)
	if err != nil {
		return gjson.Result{} // Empty result on error
	}

	// Parse back into gjson.Result so we can query later
	return gjson.ParseBytes(jsonBytes)
}
