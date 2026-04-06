package fhir

import (
	"github.com/tidwall/gjson"

	"github.com/gabrielcortez835-lab/Novellia-take-home-assessment/extractionConfig"
)

func ExtractFields(resourceJSON string, cfg extractionConfig.ExtractionConfig) map[string]gjson.Result {
	output := make(map[string]gjson.Result)
	resourceType := gjson.Get(resourceJSON, "resourceType").String()

	for field, allowed := range cfg.Fields {
		shouldExtract := false

		switch v := allowed.(type) {
		case string:
			if v == "all" {
				shouldExtract = true
			}
		case []string:
			for _, t := range v {
				if t == resourceType {
					shouldExtract = true
					break
				}
			}
		}

		if shouldExtract {
			result := gjson.Get(resourceJSON, field)
			if result.Exists() {
				output[field] = result
			}
		}
	}

	return output
}
