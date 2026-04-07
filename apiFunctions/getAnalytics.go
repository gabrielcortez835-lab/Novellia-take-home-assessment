package apiFunctions

import (
	"encoding/json"
	"log"

	"github.com/gabrielcortez835-lab/Novellia-take-home-assessment/constants"
	"github.com/gabrielcortez835-lab/Novellia-take-home-assessment/objects"
	"github.com/gabrielcortez835-lab/Novellia-take-home-assessment/sql"
)

func ApiGetAnalytics() (string, error) {

	a, err := ApiGetAnalyticsObject()

	if err != nil {
		return "", err
	}

	jsonBytes, err := json.Marshal(a) // compact JSON

	if err != nil {
		log.Fatalf("failed to marshal: %v", err)
	}

	returnString := string(jsonBytes)

	return returnString, nil
}

func ApiGetAnalyticsObject() (objects.Analytics, error) {
	a := objects.Analytics{
		RecordsByResourceType:  make(map[string]int),
		NumberOfUniqueSubjects: 0,                                  // default, optional
		ValidationErrorSummary: make([]objects.ValidationError, 0), // empty slice
		TotalEntriesPerPatient: make(map[string]int),
	}

	for resourceType := range constants.ValidResourceTypes {
		count, err := sql.GetRecordCountForResourceType(resourceType)

		if err != nil {
			return a, err
		}

		a.RecordsByResourceType[resourceType] = count
	}

	count, err := sql.GetRecordCountForSubject()

	if err != nil {
		return a, err
	}

	a.NumberOfUniqueSubjects = count

	validationErrors, err := sql.GetAllValidationErrors()

	if err != nil {
		return a, err
	}

	a.ValidationErrorSummary = validationErrors

	if err != nil {
		return a, err
	}

	subjectCountMap, err := sql.GetRecordCountPerPatient()

	if err != nil {
		return a, err
	}

	a.TotalEntriesPerPatient = subjectCountMap

	return a, nil
}
