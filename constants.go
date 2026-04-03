package main

const SqlDBName string = "FHIRDataDB"
const FHIRDataTableName string = "FHIRData"

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
