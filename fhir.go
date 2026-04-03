package main

type FHIRBase struct {
	ResourceType string `json:"resourceType"`
	ID           string `json:"id"`
}

type Coding struct {
	System  string `json:"system,omitempty"`
	Code    string `json:"code,omitempty"`
	Display string `json:"display,omitempty"`
}

type CodeableConcept struct {
	Coding []Coding `json:"coding"`
	Text   string   `json:"text,omitempty"`
}

type Quantity struct {
	Value float64 `json:"value"`
	Unit  string  `json:"unit,omitempty"`
}

type Reference struct {
	Reference string `json:"reference"`
	Display   string `json:"display,omitempty"`
}

type DosageInstruction struct {
	Text string `json:"text"`
}

type HumanName struct {
	Use    string   `json:"use,omitempty"`
	Family string   `json:"family,omitempty"`
	Given  []string `json:"given,omitempty"`
}

type ObservationComponent struct {
	Code          CodeableConcept `json:"code"`
	ValueQuantity Quantity        `json:"valueQuantity"`
}

/*
   Explicit FHIR Resource representations
*/

// Observation resource
type Observation struct {
	*FHIRBase
	Status            string                 `json:"status,omitempty"`
	Code              CodeableConcept        `json:"code"`
	Subject           Reference              `json:"subject"`
	EffectiveDateTime string                 `json:"effectiveDateTime,omitempty"`
	ValueQuantity     *Quantity              `json:"valueQuantity,omitempty"`
	Component         []ObservationComponent `json:"component,omitempty"`
}

// Condition resource
type Condition struct {
	*FHIRBase
	ClinicalStatus     *CodeableConcept `json:"clinicalStatus,omitempty"`
	VerificationStatus *CodeableConcept `json:"verificationStatus,omitempty"`
	Code               CodeableConcept  `json:"code"`
	Subject            Reference        `json:"subject"`
	OnsetDateTime      string           `json:"onsetDateTime,omitempty"`
}

// MedicationRequest resource
type MedicationRequest struct {
	*FHIRBase
	Status                    string              `json:"status"`
	Intent                    string              `json:"intent"`
	MedicationCodeableConcept CodeableConcept     `json:"medicationCodeableConcept"`
	Subject                   Reference           `json:"subject"`
	AuthoredOn                string              `json:"authoredOn,omitempty"`
	DosageInstruction         []DosageInstruction `json:"dosageInstruction,omitempty"`
}

// Procedure resource
type Procedure struct {
	*FHIRBase
	Status            string          `json:"status"`
	Code              CodeableConcept `json:"code"`
	Subject           Reference       `json:"subject"`
	PerformedDateTime string          `json:"performedDateTime,omitempty"`
}

// Patient resource
type Patient struct {
	*FHIRBase
	Name   []HumanName `json:"name"`
	Gender string      `json:"gender,omitempty"`
	Active bool        `json:"active,omitempty"`
}
