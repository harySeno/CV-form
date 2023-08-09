package models

type Applicant struct {
	ProfileCode int `json:"profileCode"`
	PersonalDetail
	WorkExp
	PastJob
	Education
	Skill
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

type PastJob struct {
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

var MockApplicantData = []Applicant{{
	ProfileCode: 12345678,
	PersonalDetail: PersonalDetail{
		WantedJobTitle: "Software Engineer",
		FirstName:      "Namaku",
		LastName:       "Ukaman",
		Email:          "ukaman.namaku@gmail.com",
		Phone:          "08008880000",
		Country:        "Indonesia",
		City:           "Jakarta",
		Address:        "Jl. Gatot Subroto",
		PostalCode:     200001,
		DrivingLicense: "1234567890123456",
		Nationality:    "Indonesia",
		PlaceOfBirth:   "Maluku",
		DateOfBirth:    "07-12-1988",
		PhotoUrl:       "/app/upload/photo/12345678.png",
	},
	WorkExp: WorkExp{
		WorkingExperience: "software engineer bla bla bla",
	},
	PastJob: PastJob{
		ID:          1,
		JobTitle:    "CEO",
		Employer:    "Toko Lapak",
		StartDate:   "01-01-2020",
		EndDate:     "01-01-2021",
		City:        "Jakarta",
		Description: "CEO",
	},
	Education: Education{
		ID:          1,
		School:      "ITB",
		Degree:      "21",
		StartDate:   "01-06-2000",
		EndDate:     "01-06-2004",
		City:        "Bandung",
		Description: "ITB",
	},
	Skill: Skill{
		ID:    1,
		Skill: "Docker",
		Level: "Expert",
	},
},
}
