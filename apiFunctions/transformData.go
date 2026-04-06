package apiFunctions

import (
	"github.com/gabrielcortez835-lab/Novellia-take-home-assessment/sql"
	"github.com/tidwall/gjson"
	"github.com/tidwall/sjson"
)

func TransformRequest(jsonString string) (string, error) {

	transformObj := gjson.Parse(jsonString)

	records, err := getRecords(transformObj)

	if err != nil {
		return "", err
	}
	transformations := transformObj.Get("transformations")

	applyTransformations := transformations.Exists() && len(transformations.Map()) > 0

	jsonArr := make([]gjson.Result, 0)

	for _, record := range records {
		if applyTransformations {
			jsonStr := "{}"
			transformations.ForEach(func(_, value gjson.Result) bool {
				action := value.Get("action").String()
				switch action {
				case "flatten":
					jsonStr = flattenField(jsonStr, record, value.Get("field").String())
				case "extract":
					jsonStr = extractField(jsonStr, record, value.Get("field").String(), value.Get("as").String())
				}

				return true
			})
		} else {
			jsonArr = append(jsonArr, record)
		}
	}

	return ResultsToJSONL(jsonArr, nil)
}

func flattenField(jsonStr string, record gjson.Result, field string) string {
	val := record.Get(field)
	val.ForEach(func(key, value gjson.Result) bool {
		jsonStr, _ = sjson.Set(jsonStr, key.String(), value)
		return true
	})

	return jsonStr
}

func extractField(jsonStr string, record gjson.Result, field string, newFieldName string) string {
	jsonStr, _ = sjson.Set(jsonStr, newFieldName, record.Get(field).Value())

	return jsonStr
}

func getRecords(transformObj gjson.Result) ([]gjson.Result, error) {
	subjectString := ""
	resourceTypeString := ""

	if subjectObj := transformObj.Get("filters.subject"); subjectObj.Exists() {
		subjectString = subjectObj.Str
	}

	if resourceTypeObj := transformObj.Get("resourceType"); resourceTypeObj.Exists() {
		resourceTypeString = resourceTypeObj.Str
	}

	records, err := sql.SqlGetRecords(subjectString, resourceTypeString)

	if err != nil {
		return nil, err
	}

	return records, nil
}
