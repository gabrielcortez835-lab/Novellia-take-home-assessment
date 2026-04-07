package apiFunctions

import "github.com/gabrielcortez835-lab/Novellia-take-home-assessment/sql"

func GetRecords(resourceTypeFilter string, subjectFilter string, fields []string) (string, error) {
	sqlResults, err := sql.GetRecords(resourceTypeFilter, subjectFilter)

	if err != nil {
		return "", err
	}

	return ResultsToJSONL(sqlResults, fields)
}

func GetRecordsById(id string, fields []string) (string, error) {

	sqlResults, err := sql.GetRecordsById(id)

	if err != nil {
		return "", err
	}

	return ResultsToJSONL(sqlResults, fields)
}
