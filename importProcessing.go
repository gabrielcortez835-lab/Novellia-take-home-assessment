package main

import (
	"encoding/json" // for JSON parsing
	"fmt"           // for printing and formatting
)

func ProcessImportedJson(obj map[string]interface{}) ([]string, error) {

	resourceType := obj["resourceType"].(string)

	constructor, found := resourceFactory[resourceType]
	if !found {
		return nil, fmt.Errorf("unknown resourceType: %s", resourceType)
	}

	resource := constructor()

	jsonBytes, err := json.Marshal(obj)
	if err != nil {
		return nil, err
	}

	if err := json.Unmarshal(jsonBytes, resource); err != nil {
		return nil, err
	}

	//Load into SQL
	var patientID string
	var id string

	switch resourceType {
	case "Condition":
		cond := resource.(*Condition)

		if dataQualityWarnings := validateCondition(*cond); len(dataQualityWarnings) > 0 {
			return dataQualityWarnings, nil
		}

		patientID = cond.Subject.Reference
		id = cond.ID
	case "Observation":
		obs := resource.(*Observation)

		if dataQualityWarnings := validateObservation(*obs); len(dataQualityWarnings) > 0 {
			return dataQualityWarnings, nil
		}

		patientID = obs.Subject.Reference
		id = obs.ID
	case "MedicationRequest":
		medReq := resource.(*MedicationRequest)

		if dataQualityWarnings := validateMedicationRequest(*medReq); len(dataQualityWarnings) > 0 {
			return dataQualityWarnings, nil
		}

		patientID = medReq.Subject.Reference
		id = medReq.ID
	case "Procedure":
		proc := resource.(*Procedure)

		if dataQualityWarnings := validateProcedure(*proc); len(dataQualityWarnings) > 0 {
			return dataQualityWarnings, nil
		}

		patientID = proc.Subject.Reference
		id = proc.ID
	case "Patient":
		patient := resource.(*Patient)

		if dataQualityWarnings := validatePatient(*patient); len(dataQualityWarnings) > 0 {
			return dataQualityWarnings, nil
		}

		patientID = patient.ID
		id = patient.ID
	}

	db, err := Connect(SqlDBName)

	if err != nil {
		return nil, err
	}

	defer db.Close()

	// Convert bytes to string
	jsonString := string(jsonBytes)

	_, err = InsertResource(db, id, resourceType, patientID, jsonString)

	if err != nil {
		return nil, err
	}

	return nil, nil
}

type ResourceConstructor func() interface{}

var resourceFactory = map[string]ResourceConstructor{
	"Observation":       func() interface{} { return &Observation{} },
	"Condition":         func() interface{} { return &Condition{} },
	"MedicationRequest": func() interface{} { return &MedicationRequest{} },
	"Procedure":         func() interface{} { return &Procedure{} },
	"Patient":           func() interface{} { return &Patient{} },
}
