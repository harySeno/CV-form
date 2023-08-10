package models

import "github.com/jinzhu/gorm"

type Applicant struct {
	gorm.Model
	ProfileCode int `json:"profileCode"`
	PersonalDetail
	WorkExp
	Employment []Employment `json:"employment"`
	Education  []Education  `json:"education"`
	Skill      []Skill      `json:"skill"`
}

type PersonalDetail struct {
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

type WorkExp struct {
	WorkingExperience string `json:"workingExperience"`
}

type Employment struct {
	ID          int    `json:"id"`
	JobTitle    string `json:"jobTitle"`
	Employer    string `json:"employer"`
	StartDate   string `json:"startDate"`
	EndDate     string `json:"endDate"`
	City        string `json:"city"`
	Description string `json:"description"`
}

type Education struct {
	ID          int    `json:"id"`
	School      string `json:"school"`
	Degree      string `json:"degree"`
	StartDate   string `json:"startDate"`
	EndDate     string `json:"endDate"`
	City        string `json:"city"`
	Description string `json:"description"`
}

type Skill struct {
	ID    int    `json:"id"`
	Skill string `json:"skill"`
	Level string `json:"level"`
}

func (Applicant) TableName() string {
	return "applicants"
}

func (PersonalDetail) TableName() string {
	return "personaldetails"
}

func (WorkExp) TableName() string {
	return "workexps"
}

func (Employment) TableName() string {
	return "employments"
}

func (Education) TableName() string {
	return "educations"
}

func (Skill) TableName() string {
	return "skills"
}
