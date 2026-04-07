package constants

const SqlDBName = "FHIRDataDB"
const FHIRDataTableName = "FHIRData"
const ValidationErrorTableName = "ValidationErrorTable"

var ValidStatus = map[string]bool{
	"active":    true,
	"completed": true,
	"final":     true,
}

var ValidResourceTypes = map[string]bool{
	"Observation":       true,
	"Condition":         true,
	"MedicationRequest": true,
	"Procedure":         true,
	"Patient":           true,
}

const (
	ApiPostImportPath     = "/import"
	ApiGetRecordsByIdPath = "/records/:id"
	ApiGetRecordsPath     = "/records"
	ApiPostTransformPath  = "/transform"
	ApiGetAnalytics       = "/analytics"
)

const ExtractionConfigFileName = "extractionConfig/extractionConfig.json"

var ResourceTypeEnum = struct {
	Patient           string
	Observation       string
	MedicationRequest string
	Procedure         string
	Condition         string
}{
	Patient:           "Patient",
	Observation:       "Observation",
	MedicationRequest: "MedicationRequest",
	Procedure:         "Procedure",
	Condition:         "Condition",
}
