package restapi

import (
	"CV-form/pkg/models"
	"errors"
	"github.com/emicklei/go-restful/v3"
	"net/http"
	"strconv"
)

var candidate []models.Applicant

// InitializeMockData initializes the mock data for testing
func InitializeMockData() {
	candidate = models.MockApplicantData
}

func GetAllData(_ *restful.Request, res *restful.Response) {
	err := res.WriteEntity(candidate)
	if err != nil {
		return
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
			result := struct {
				ProfileCode int `json:"profileCode"`
				models.PersonalDetail
			}{
				ProfileCode:    applicant.ProfileCode,
				PersonalDetail: applicant.PersonalDetail,
			}
			err := response.WriteEntity(result)
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

	// Generate a new profile code and assign it
	newProfileCode := len(candidate) + 1
	addRequest.ProfileCode = newProfileCode

	// Append the new applicant to the candidate list
	candidate = append(candidate, *addRequest)

	// Create a response JSON with the new profile code
	result := struct {
		ProfileCode int `json:"profileCode"`
	}{
		ProfileCode: newProfileCode,
	}

	err = res.WriteHeaderAndEntity(http.StatusCreated, result)
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

	// Read the updated data from the request
	updateRequest := &models.Applicant{}
	err = request.ReadEntity(updateRequest)
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

	candidate[indexToUpdate] = models.Applicant{
		ProfileCode:    code,
		PersonalDetail: updateRequest.PersonalDetail,
	}

	err = response.WriteHeaderAndEntity(http.StatusOK, candidate[indexToUpdate])
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
			resp := struct {
				WorkingExperience string `json:"workingExperience"`
			}{
				WorkingExperience: applicant.WorkExp.WorkingExperience,
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

	candidate[indexToUpdate].WorkExp.WorkingExperience = updateRequest.WorkingExperience
	result := struct {
		ProfileCode string `json:"profileCode"` // is it intended to show profileCode as the field name instead of workingExperience?
	}{
		ProfileCode: updateRequest.WorkingExperience,
	}

	err = response.WriteHeaderAndEntity(http.StatusOK, result)
	if err != nil {
		return
	}
}

// GetEmploymentByCode handles GET requests to retrieve an applicant employment by profile code
func GetEmploymentByCode(request *restful.Request, response *restful.Response) {
	candidateCode := request.PathParameter("code")
	code, err := strconv.Atoi(candidateCode)
	if err != nil {
		err = response.WriteError(http.StatusBadRequest, err)
		return
	}

	for _, applicant := range candidate {
		if applicant.ProfileCode == code {
			resp := struct {
				Employment []models.Employment `json:"employment"`
			}{
				Employment: applicant.Employment,
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

// AddEmployment handles POST requests to add an applicant employment by profile code
func AddEmployment(request *restful.Request, response *restful.Response) {
	candidateCode := request.PathParameter("code")
	code, err := strconv.Atoi(candidateCode)
	if err != nil {
		err = response.WriteError(http.StatusBadRequest, err)
		return
	}

	employment := models.Employment{}
	err = request.ReadEntity(&employment)
	if err != nil {
		err = response.WriteError(http.StatusBadRequest, err)
		return
	}

	// Find the applicant with the given code
	var targetApplicant *models.Applicant
	for i := range candidate {
		if candidate[i].ProfileCode == code {
			targetApplicant = &candidate[i]
			break
		}
	}

	if targetApplicant == nil {
		err = response.WriteError(http.StatusNotFound, errors.New("applicant not found"))
		if err != nil {
			return
		}
		return
	}

	// Generate a new ID for the new employment
	newID := len(targetApplicant.Employment) + 1
	employment.ID = newID

	// Append the new employment to the applicant's employment list
	targetApplicant.Employment = append(targetApplicant.Employment, employment)

	// Return the added employment data
	err = response.WriteHeaderAndEntity(http.StatusCreated, employment)
	if err != nil {
		return
	}
}
