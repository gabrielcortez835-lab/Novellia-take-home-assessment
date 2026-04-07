package apiFunctions

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/tidwall/gjson"
)

func ResultsToJSONL(results []gjson.Result, fields []string) (string, error) {
	fmt.Println(len(fields))
	var jsonlBuilder strings.Builder
	if len(fields) > 0 {
		for _, record := range results {
			obj := make(map[string]interface{})
			for _, field := range fields {
				val := record.Get(field)
				if val.Exists() {
					obj[field] = val.Value()
				}
			}
			lineBytes, err := json.Marshal(obj)
			if err != nil {
				return "", fmt.Errorf("failed to marshal JSON: %w", err)
			}
			jsonlBuilder.Write(lineBytes)
			jsonlBuilder.WriteString("\n")
		}
	} else {
		for _, record := range results {
			jsonlBuilder.WriteString(record.Raw)
			jsonlBuilder.WriteString("\n")
		}
	}

	return jsonlBuilder.String(), nil
}

func StructSliceToJSONStringSlice[T any](items []T) ([]string, error) {
	out := make([]string, 0, len(items))
	for _, v := range items {
		b, err := json.Marshal(v)
		if err != nil {
			return nil, err
		}
		out = append(out, string(b))
	}
	return out, nil
}
