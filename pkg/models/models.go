package models

type Applicant struct {
	ProfileCode    int    `json:"profileCode"`
	WantedJobTitle string `json:"wantedJobTitle"`
	FirstName      string `json:"firstName"`
	LastName       string `json:"lastName"`
	Email          string `json:"email"`
	Phone          string `json:"phone"`
	Country        string `json:"country"`
	City           string `json:"city"`
	Address        string `json:"address"`
	PostalCode     int    `json:"postalCode"`
	DrivingLicense string `json:"drivingLicense"`
	Nationality    string `json:"nationality"`
	PlaceOfBirth   string `json:"placeOfBirth"`
	DateOfBirth    string `json:"dateOfBirth"`
	PhotoUrl       string `json:"photoUrl"`
}
