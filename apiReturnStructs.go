package main

type ImportReturn struct {
	TotalLinesProcessed         int                 `json:"totalLinesProcessed"`
	RecordsImportedSuccessfully int                 `json:"recordsImportedSuccessfully"`
	ValidationErrors            []string            `json:"validationErrors"`
	DataQualityWarnings         map[string][]string `json:"dataQualityWarnings"`
	Statistics                  Statistics          `json:"statistics"`
}

type Statistics struct {
	RecordsByType  map[string]int `json:"recordsByType"`
	UniquePatients int            `json:"uniquePatients"`
}
