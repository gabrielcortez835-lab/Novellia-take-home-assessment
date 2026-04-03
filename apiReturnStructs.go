package main

type ImportReturn struct {
	totalLinesProcessed         int
	recordsImportedSuccessfully int
	validationErrors            []string
	dataQualityWarnings         map[string][]string
	statistics                  Statistics
}

type Statistics struct {
	RecordsByType  map[string]int
	uniquePatients int
}
