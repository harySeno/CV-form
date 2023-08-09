package restapi

import (
	"CV-form/pkg/models"
	"encoding/base64"
	"errors"
	"fmt"
	"github.com/emicklei/go-restful/v3"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
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

// UploadPhoto handles PUT requests to upload applicant's photo
func UploadPhoto(request *restful.Request, response *restful.Response) {
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
	updateRequest := struct {
		Image string `json:"base64img"`
	}{}

	err = request.ReadEntity(&updateRequest)
	if err != nil {
		err = response.WriteError(http.StatusBadRequest, err)
		return
	}

	// Save the base64 image to a file
	b64data := updateRequest.Image[strings.IndexByte(updateRequest.Image, ',')+1:]
	imageFile, err := base64.StdEncoding.DecodeString(b64data)
	if err != nil {
		err = response.WriteError(http.StatusInternalServerError, err)
		return
	}

	filename := fmt.Sprintf("%d.png", candidate[indexToUpdate].ProfileCode)
	imagePath := "/app/upload/photo/"

	// Write the image data to the file
	err = ioutil.WriteFile(imagePath+filename, imageFile, 0644)
	if err != nil {
		err = response.WriteError(http.StatusInternalServerError, err)
		return
	}

	// Update the PhotoUrl field with the new path
	candidate[indexToUpdate].PhotoUrl = imagePath + filename
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
			result := struct {
				WorkingExperience string `json:"workingExperience"`
			}{
				WorkingExperience: applicant.WorkExp.WorkingExperience,
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
			result := struct {
				Employment []models.Employment `json:"employment"`
			}{
				Employment: applicant.Employment,
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

	result := struct {
		ProfileCode int `json:"profileCode"`
		ID          int `json:"id"`
	}{
		ProfileCode: code,
		ID:          newID,
	}

	err = response.WriteHeaderAndEntity(http.StatusCreated, result)
	if err != nil {
		return
	}
}

// DeleteEmployment handles DELETE requests to remove an applicant employment by profile code
func DeleteEmployment(request *restful.Request, response *restful.Response) {
	candidateCode := request.PathParameter("code")
	code, err := strconv.Atoi(candidateCode)
	if err != nil {
		err = response.WriteError(http.StatusBadRequest, err)
		return
	}

	// Parse the 'id' query parameter
	employmentID := request.QueryParameter("id")
	id, err := strconv.Atoi(employmentID)
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

	// Find the index of the employment with the specified 'id'
	indexToRemove := -1
	for i := range targetApplicant.Employment {
		if targetApplicant.Employment[i].ID == id {
			indexToRemove = i
			break
		}
	}

	if indexToRemove == -1 {
		err = response.WriteError(http.StatusNotFound, errors.New("employment not found"))
		if err != nil {
			return
		}
		return
	}

	// Remove the employment from the applicant's Employment slice
	targetApplicant.Employment = append(targetApplicant.Employment[:indexToRemove], targetApplicant.Employment[indexToRemove+1:]...)

	result := struct {
		ProfileCode int `json:"profileCode"`
	}{
		ProfileCode: code,
	}

	err = response.WriteEntity(result)
	if err != nil {
		return
	}
}

// GetEducationByCode handles GET requests to retrieve an applicant education by profile code
func GetEducationByCode(request *restful.Request, response *restful.Response) {
	candidateCode := request.PathParameter("code")
	code, err := strconv.Atoi(candidateCode)
	if err != nil {
		err = response.WriteError(http.StatusBadRequest, err)
		return
	}

	for _, applicant := range candidate {
		if applicant.ProfileCode == code {
			result := struct {
				Education []models.Education `json:"education"`
			}{
				Education: applicant.Education,
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

// AddEducation handles POST requests to add an applicant education by profile code
func AddEducation(request *restful.Request, response *restful.Response) {
	candidateCode := request.PathParameter("code")
	code, err := strconv.Atoi(candidateCode)
	if err != nil {
		err = response.WriteError(http.StatusBadRequest, err)
		return
	}

	education := models.Education{}
	err = request.ReadEntity(&education)
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

	// Generate a new ID for the new education
	newID := len(targetApplicant.Education) + 1
	education.ID = newID

	// Append the new education to the applicant's education list
	targetApplicant.Education = append(targetApplicant.Education, education)

	result := struct {
		ProfileCode int `json:"profileCode"`
		ID          int `json:"id"`
	}{
		ProfileCode: code,
		ID:          newID,
	}

	err = response.WriteHeaderAndEntity(http.StatusCreated, result)
	if err != nil {
		return
	}
}

// DeleteEducation handles DELETE requests to remove an applicant education by profile code
func DeleteEducation(request *restful.Request, response *restful.Response) {
	candidateCode := request.PathParameter("code")
	code, err := strconv.Atoi(candidateCode)
	if err != nil {
		err = response.WriteError(http.StatusBadRequest, err)
		return
	}

	// Parse the 'id' query parameter
	educationID := request.QueryParameter("id")
	id, err := strconv.Atoi(educationID)
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

	// Find the index of the education with the specified 'id'
	indexToRemove := -1
	for i := range targetApplicant.Education {
		if targetApplicant.Education[i].ID == id {
			indexToRemove = i
			break
		}
	}

	if indexToRemove == -1 {
		err = response.WriteError(http.StatusNotFound, errors.New("education not found"))
		if err != nil {
			return
		}
		return
	}

	// Remove the education from the applicant's Education slice
	targetApplicant.Education = append(targetApplicant.Education[:indexToRemove], targetApplicant.Education[indexToRemove+1:]...)

	result := struct {
		ProfileCode int `json:"profileCode"`
	}{
		ProfileCode: code,
	}

	err = response.WriteEntity(result)
	if err != nil {
		return
	}
}
