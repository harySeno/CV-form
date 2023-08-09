package restapi

import (
	"CV-form/pkg/models"
	"errors"
	"github.com/emicklei/go-restful/v3"
	"net/http"
	"strconv"
)

// candidate holds the mock data for applicants
var candidate []models.Applicant

// InitializeMockData initializes the mock data for testing
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

// GetByCode handles GET requests to retrieve an applicant by profile code
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

// AddProfile handles POST requests to add a new applicant profile
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

	// Create a new applicant with data from the request
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

	// Generate a new profile code and assign it
	newProfileCode := len(candidate) + 1
	applicant.ProfileCode = newProfileCode

	// Append the new applicant to the candidate list
	candidate = append(candidate, applicant)

	// Create a response JSON with the new profile code
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

// UpdateProfile handles PUT requests to update an existing applicant profile by profile code
func UpdateProfile(request *restful.Request, response *restful.Response) {
	candidateCode := request.PathParameter("code")
	code, err := strconv.Atoi(candidateCode)
	if err != nil {
		err = response.WriteError(http.StatusBadRequest, err)
		return
	}

	// Find the index of the profile with the given code
	indexToUpdate := -1
	for i, applicant := range candidate {
		if applicant.ProfileCode == code {
			indexToUpdate = i
			break
		}
	}

	if indexToUpdate == -1 {
		err = response.WriteError(http.StatusNotFound, errors.New("profile not found"))
		if err != nil {
			return
		}
		return
	}

	// Read the updated data from the request
	updateRequest := &models.Applicant{}
	err = request.ReadEntity(updateRequest)
	if err != nil {
		err = response.WriteError(http.StatusBadRequest, err)
		return
	}

	// Update the profile data
	candidate[indexToUpdate] = models.Applicant{
		ProfileCode:    code,
		WantedJobTitle: updateRequest.WantedJobTitle,
		FirstName:      updateRequest.FirstName,
		LastName:       updateRequest.LastName,
		Email:          updateRequest.Email,
		Phone:          updateRequest.Phone,
		Country:        updateRequest.Country,
		City:           updateRequest.City,
		Address:        updateRequest.Address,
		PostalCode:     updateRequest.PostalCode,
		DrivingLicense: updateRequest.DrivingLicense,
		Nationality:    updateRequest.Nationality,
		PlaceOfBirth:   updateRequest.PlaceOfBirth,
		DateOfBirth:    updateRequest.DateOfBirth,
		PhotoUrl:       updateRequest.PhotoUrl,
	}

	err = response.WriteHeaderAndEntity(http.StatusOK, candidate)
	if err != nil {
		return
	}
}

// GetExpByCode handles GET requests to retrieve an applicant working experience by profile code
func GetExpByCode(request *restful.Request, response *restful.Response) {
	candidateCode := request.PathParameter("code")
	code, err := strconv.Atoi(candidateCode)
	if err != nil {
		err = response.WriteError(http.StatusBadRequest, err)
		return
	}
	for _, applicant := range candidate {
		if applicant.ProfileCode == code {
			// Create a response JSON with the new workingExperience
			workingExperience := "Software Engineer with bla bla bla experience."
			resp := struct {
				WorkingExperience string `json:"workingExperience"`
			}{
				WorkingExperience: workingExperience,
			}
			err := response.WriteEntity(resp)
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

// UpdateExperience handles PUT requests to update an existing applicant working experience by profile code
func UpdateExperience(request *restful.Request, response *restful.Response) {
	candidateCode := request.PathParameter("code")
	code, err := strconv.Atoi(candidateCode)
	if err != nil {
		err = response.WriteError(http.StatusBadRequest, err)
		return
	}

	updateRequest := struct {
		WorkingExperience string `json:"workingExperience"`
	}{}

	err = request.ReadEntity(&updateRequest)
	if err != nil {
		err = response.WriteError(http.StatusBadRequest, err)
		return
	}

	// Find the index of the profile with the given code
	indexToUpdate := -1
	for i, applicant := range candidate {
		if applicant.ProfileCode == code {
			indexToUpdate = i
			break
		}
	}

	if indexToUpdate == -1 {
		err = response.WriteError(http.StatusNotFound, errors.New("profile not found"))
		if err != nil {
			return
		}
		return
	}

	candidate[indexToUpdate].WorkingExperience = updateRequest.WorkingExperience

	err = response.WriteHeaderAndEntity(http.StatusOK, candidate[indexToUpdate])
	if err != nil {
		return
	}
}
