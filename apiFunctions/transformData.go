package apiFunctions

import (
	"fmt"
	"regexp"

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

	applyTransformations := transformations.Exists() && len(transformObj.Array()) > 0

	jsonArr := make([]gjson.Result, 0)

	fmt.Printf("records: %d \n", len(records))

	fmt.Println(records)

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
			jsonArr = append(jsonArr, gjson.Parse(jsonStr))
		} else {
			jsonArr = append(jsonArr, record)
		}
	}

	return ResultsToJSONL(jsonArr, nil)
}

var arrayIndexPattern = regexp.MustCompile(`\[(\d+)\]`)

func flattenField(jsonStr string, record gjson.Result, field string) string {
	// Convert array notation [0] → .0 for gjson
	gjsonField := arrayIndexPattern.ReplaceAllString(field, `.$1`)

	target := record.Get(gjsonField)
	if !target.Exists() {
		fmt.Println("field not found:", gjsonField)
		return jsonStr
	}

	// Flatten object keys
	if target.IsObject() {
		target.ForEach(func(key, value gjson.Result) bool {
			jsonStr, _ = sjson.Set(jsonStr, key.String(), value.Value())
			return true
		})
	}

	// Flatten each element in an array (rare case you want them all)
	if target.IsArray() {
		for _, el := range target.Array() {
			if el.IsObject() {
				el.ForEach(func(key, value gjson.Result) bool {
					jsonStr, _ = sjson.Set(jsonStr, key.String(), value.Value())
					return true
				})
			}
		}
	}

	return jsonStr
}

func extractField(jsonStr string, record gjson.Result, field string, newFieldName string) string {
	fmt.Println("extract:", field)

	// Convert array notation for gjson
	gjsonField := arrayIndexPattern.ReplaceAllString(field, `.$1`)

	val := record.Get(gjsonField)
	if !val.Exists() {
		fmt.Println("field not found:", gjsonField)
		return jsonStr
	}

	// Assign value to new field name
	jsonStr, _ = sjson.Set(jsonStr, newFieldName, val.Value())
	fmt.Println(jsonStr)
	return jsonStr
}

func getRecords(transformObj gjson.Result) ([]gjson.Result, error) {
	subjectString := ""

	// Check subject filter
	if subjectObj := transformObj.Get("filters.subject"); subjectObj.Exists() {
		subjectString = subjectObj.Str
	}

	records := []gjson.Result{}

	// Loop over resourceTypes array
	resourceTypes := transformObj.Get("resourceTypes").Array()
	if len(resourceTypes) == 0 {
		// fallback if single resourceType provided instead of array
		if resourceTypeObj := transformObj.Get("resourceType"); resourceTypeObj.Exists() {
			resourceTypes = []gjson.Result{resourceTypeObj}
		}
	}

	// Iterate through each resourceType
	for _, rt := range resourceTypes {
		rtStr := rt.Str
		rs, err := sql.GetRecords(rtStr, subjectString)
		if err != nil {
			return nil, err // stop on first error
		}
		records = append(records, rs...)
	}

	return records, nil
}
