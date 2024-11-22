package models

import "time"

type Kev struct {
	CveID             string    `json:"cveID" xml:"cveID"`
	VendorProject     string    `json:"vendorProject" xml:"vendorProject"`
	Product           string    `json:"product" xml:"product"`
	VulnerabilityName string    `json:"vulnerabilityName" xml:"vulnerabilityName"`
	ShortDescription  string    `json:"shortDescription" xml:"shortDescription"`
	RequiredAction    string    `json:"requiredAction" xml:"requiredAction"`
	DateAdded         time.Time `json:"dateAdded" xml:"dateAdded"`
	Notes             string    `json:"notes" xml:"notes"`
	Cews              any       `json:"cews" xml:"cews"`
}
