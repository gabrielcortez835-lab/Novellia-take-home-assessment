package objects

type ValidationError struct {
	EntryID         string   `json:"EntryID"`
	ValidationError []string `json:"ValidationError"`
}

type Analytics struct {
	RecordsByResourceType  map[string]int    `json:"recordsByResourceType"`
	NumberOfUniqueSubjects int               `json:"numberOfUniqueSubjects"`
	ValidationErrorSummary []ValidationError `json:"validationErrorSummary"`
	TotalEntriesPerPatient map[string]int    `json:"totalEntriesPerPatient"`
}
