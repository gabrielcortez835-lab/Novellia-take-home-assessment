package constants

const SqlDBName = "FHIRDataDB"
const FHIRDataTableName = "FHIRData"

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

const ApiPostImportPath = "/import"
const ApiGetRecordsByIdPath = "/records/:id"
const ApiGetRecordsPath = "/records"
const ApiPostTransformPath = "/transform"

const ExtractionConfigFileName = "extactionConfig.json"


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
