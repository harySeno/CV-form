package restapi

import (
	"CV-form/pkg/models"
	"github.com/emicklei/go-restful/v3"
	"net/http"
	"strconv"
)

var candidate []models.Applicant

func InitializeMockData() {
	candidate = []models.Applicant{{
		ProfileCode:    12345678,
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
	}
}

func GetByCode(request *restful.Request, response *restful.Response) {
	candidateCode := request.PathParameter("code")
	code, err := strconv.Atoi(candidateCode)
	if err != nil {
		err = response.WriteError(http.StatusBadRequest, err)
		return
	}
	for _, applicant := range candidate {
		if applicant.ProfileCode == code {
			err := response.WriteEntity(applicant)
			if err != nil {
				return
			}
			return
		}
	}
	err = response.WriteError(http.StatusNotFound, err)
	if err != nil {
		return
	}
}

func AddProfile(req *restful.Request, res *restful.Response) {
	addRequest := &models.Applicant{}

	err := req.ReadEntity(addRequest)
	if err != nil {
		err := res.WriteErrorString(http.StatusBadRequest, "Invalid payload format")
		if err != nil {
			return
		}
		return
	}

	applicant := models.Applicant{
		WantedJobTitle: addRequest.WantedJobTitle,
		FirstName:      addRequest.FirstName,
		LastName:       addRequest.LastName,
		Email:          addRequest.Email,
		Phone:          addRequest.Phone,
		Country:        addRequest.Country,
		City:           addRequest.City,
		Address:        addRequest.Address,
		PostalCode:     addRequest.PostalCode,
		DrivingLicense: addRequest.DrivingLicense,
		Nationality:    addRequest.Nationality,
		PlaceOfBirth:   addRequest.PlaceOfBirth,
		DateOfBirth:    addRequest.DateOfBirth,
		PhotoUrl:       addRequest.PhotoUrl,
	}

	newProfileCode := len(candidate) + 1
	applicant.ProfileCode = newProfileCode
	candidate = append(candidate, applicant)

	response := struct {
		ProfileCode int `json:"profileCode"`
	}{
		ProfileCode: newProfileCode,
	}
	err = res.WriteHeaderAndEntity(http.StatusCreated, response)
	if err != nil {
		return
	}

}
