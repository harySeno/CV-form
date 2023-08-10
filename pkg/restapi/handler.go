package restapi

import (
	"CV-form/pkg/database"
	"CV-form/pkg/models"
	"encoding/base64"
	"errors"
	"fmt"
	"github.com/emicklei/go-restful/v3"
	"github.com/jinzhu/gorm"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"strings"
)

var candidate []models.Applicant

func GetAllData(_ *restful.Request, response *restful.Response) {
	// Use the DB instance from database package to retrieve data
	var candidates []models.Applicant
	err := database.DB.Find(&candidates).Error
	if err != nil {
		err = response.WriteError(http.StatusNotFound, err)
		return
	}

	err = response.WriteEntity(candidate)
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

	var applicant models.Applicant
	err = database.DB.Where("profile_code = ?", code).First(&applicant).Error
	if err != nil {
		err = response.WriteError(http.StatusNotFound, err)
		return
	}

	result := struct {
		ProfileCode int `json:"profileCode"`
		models.PersonalDetail
	}{
		ProfileCode:    applicant.ProfileCode,
		PersonalDetail: applicant.PersonalDetail,
	}

	err = response.WriteEntity(result)
	if err != nil {
		return
	}
}

// AddProfile handles POST requests to add a new applicant profile
func AddProfile(request *restful.Request, response *restful.Response) {
	addRequest := &models.Applicant{}

	err := request.ReadEntity(addRequest)
	if err != nil {
		err := response.WriteErrorString(http.StatusBadRequest, "Invalid payload format")
		if err != nil {
			return
		}
		return
	}

	// Generate a new profile code and assign it
	newProfileCode := len(candidate) + 1
	addRequest.ProfileCode = newProfileCode

	// Create a new applicant in the database
	err = database.DB.Create(addRequest).Error
	if err != nil {
		err = response.WriteError(http.StatusInternalServerError, err)
		return
	}

	result := struct {
		ProfileCode int `json:"profileCode"`
	}{
		ProfileCode: addRequest.ProfileCode,
	}

	err = response.WriteHeaderAndEntity(http.StatusCreated, result)
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

	// Find the applicant in the database
	var applicant models.Applicant
	err = database.DB.Where("profile_code = ?", code).First(&applicant).Error
	if err != nil {
		if gorm.IsRecordNotFoundError(err) {
			err = response.WriteError(http.StatusNotFound, errors.New("profile not found"))
		} else {
			err = response.WriteError(http.StatusInternalServerError, err)
		}
		if err != nil {
			return
		}
		return
	}

	// Update the applicant's personal detail
	applicant.PersonalDetail = updateRequest.PersonalDetail

	// Save the updated applicant to the database
	err = database.DB.Save(&applicant).Error
	if err != nil {
		err = response.WriteError(http.StatusInternalServerError, err)
		return
	}

	result := struct {
		ProfileCode int `json:"profileCode"`
	}{
		ProfileCode: code,
	}
	err = response.WriteHeaderAndEntity(http.StatusOK, result)
	if err != nil {
		err = response.WriteError(http.StatusInternalServerError, err)
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

	// Find the applicant in the database
	var applicant models.Applicant
	err = database.DB.Where("profile_code = ?", code).First(&applicant).Error
	if err != nil {
		if gorm.IsRecordNotFoundError(err) {
			err = response.WriteError(http.StatusNotFound, errors.New("profile not found"))
		} else {
			err = response.WriteError(http.StatusInternalServerError, err)
		}
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

	filename := fmt.Sprintf("%d.png", applicant.ProfileCode)
	imagePath := "/home/seno/Desktop"

	// Write the image data to the file
	err = ioutil.WriteFile(imagePath+filename, imageFile, 0644)
	if err != nil {
		err = response.WriteError(http.StatusInternalServerError, err)
		return
	}

	// Update the PhotoUrl field with the new path
	applicant.PhotoUrl = imagePath + filename

	// Save the updated applicant to the database
	err = database.DB.Save(&applicant).Error
	if err != nil {
		err = response.WriteError(http.StatusInternalServerError, err)
		return
	}

	result := struct {
		ProfileCode int    `json:"profileCode"`
		PhotoUrl    string `json:"photoUrl"`
	}{
		ProfileCode: code,
		PhotoUrl:    applicant.PhotoUrl,
	}
	err = response.WriteHeaderAndEntity(http.StatusOK, result)
	if err != nil {
		err = response.WriteError(http.StatusInternalServerError, err)
		return
	}
}

// DownloadPhoto handles GET requests to download an applicant's photo by profile code
func DownloadPhoto(request *restful.Request, response *restful.Response) {
	candidateCode := request.PathParameter("code")
	code, err := strconv.Atoi(candidateCode)
	if err != nil {
		err = response.WriteError(http.StatusBadRequest, err)
		return
	}

	// Find the applicant in the database
	var applicant models.Applicant
	err = database.DB.Where("profile_code = ?", code).First(&applicant).Error
	if err != nil {
		if gorm.IsRecordNotFoundError(err) {
			err = response.WriteError(http.StatusNotFound, errors.New("applicant not found"))
		} else {
			err = response.WriteError(http.StatusInternalServerError, err)
		}
		if err != nil {
			return
		}
		return
	}

	// Read the image file
	imagePath := applicant.PhotoUrl
	imageFile, err := ioutil.ReadFile(imagePath)
	if err != nil {
		err = response.WriteError(http.StatusInternalServerError, err)
		return
	}

	response.Header().Set("Content-Type", "image/png")
	_, err = response.Write(imageFile)
	if err != nil {
		_ = response.WriteError(http.StatusInternalServerError, err)
		return
	}
}

// DeletePhoto handles DELETE requests to delete an applicant's photo by profile code
func DeletePhoto(request *restful.Request, response *restful.Response) {
	candidateCode := request.PathParameter("code")
	code, err := strconv.Atoi(candidateCode)
	if err != nil {
		err = response.WriteError(http.StatusBadRequest, err)
		return
	}

	// Find the applicant in the database
	var applicant models.Applicant
	err = database.DB.Where("profile_code = ?", code).First(&applicant).Error
	if err != nil {
		if gorm.IsRecordNotFoundError(err) {
			err = response.WriteError(http.StatusNotFound, errors.New("applicant not found"))
		} else {
			err = response.WriteError(http.StatusInternalServerError, err)
		}
		if err != nil {
			return
		}
		return
	}

	// Delete the image file
	imagePath := applicant.PhotoUrl
	err = os.Remove(imagePath)
	if err != nil {
		err = response.WriteError(http.StatusInternalServerError, err)
		return
	}

	// Clear the PhotoUrl field in the database
	applicant.PhotoUrl = ""
	err = database.DB.Save(&applicant).Error
	if err != nil {
		err = response.WriteError(http.StatusInternalServerError, err)
		return
	}

	result := struct {
		ProfileCode int    `json:"profileCode"`
		Message     string `json:"message"`
	}{
		ProfileCode: code,
		Message:     "Photo deleted successfully",
	}

	err = response.WriteEntity(result)
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

	// Find the applicant in the database
	var applicant models.Applicant
	err = database.DB.Where("profile_code = ?", code).First(&applicant).Error
	if err != nil {
		if gorm.IsRecordNotFoundError(err) {
			err = response.WriteError(http.StatusNotFound, errors.New("applicant not found"))
		} else {
			err = response.WriteError(http.StatusInternalServerError, err)
		}
		if err != nil {
			return
		}
		return
	}

	result := struct {
		WorkingExperience string `json:"workingExperience"`
	}{
		WorkingExperience: applicant.WorkExp.WorkingExperience,
	}

	err = response.WriteEntity(result)
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

	// Find the applicant in the database
	var applicant models.Applicant
	err = database.DB.Where("profile_code = ?", code).First(&applicant).Error
	if err != nil {
		if gorm.IsRecordNotFoundError(err) {
			err = response.WriteError(http.StatusNotFound, errors.New("applicant not found"))
		} else {
			err = response.WriteError(http.StatusInternalServerError, err)
		}
		if err != nil {
			return
		}
		return
	}

	// Update the working experience
	applicant.WorkExp.WorkingExperience = updateRequest.WorkingExperience
	err = database.DB.Save(&applicant).Error
	if err != nil {
		err = response.WriteError(http.StatusInternalServerError, err)
		return
	}

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

	// Find the applicant in the database
	var applicant models.Applicant
	err = database.DB.Preload("Employment").Where("profile_code = ?", code).First(&applicant).Error
	if err != nil {
		if gorm.IsRecordNotFoundError(err) {
			err = response.WriteError(http.StatusNotFound, errors.New("applicant not found"))
		} else {
			err = response.WriteError(http.StatusInternalServerError, err)
		}
		if err != nil {
			return
		}
		return
	}

	result := struct {
		Employment []models.Employment `json:"employment"`
	}{
		Employment: applicant.Employment,
	}
	err = response.WriteEntity(result)
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

	// Find the applicant in the database
	var applicant models.Applicant
	err = database.DB.Where("profile_code = ?", code).First(&applicant).Error
	if err != nil {
		if gorm.IsRecordNotFoundError(err) {
			err = response.WriteError(http.StatusNotFound, errors.New("applicant not found"))
		} else {
			err = response.WriteError(http.StatusInternalServerError, err)
		}
		if err != nil {
			return
		}
		return
	}

	// Generate a new ID for the new employment
	newID := len(applicant.Employment) + 1
	employment.ID = newID

	// Append the new employment to the applicant's employment list
	applicant.Employment = append(applicant.Employment, employment)

	err = database.DB.Save(&applicant).Error
	if err != nil {
		err = response.WriteError(http.StatusInternalServerError, err)
		return
	}

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

	// Find the applicant in the database
	var applicant models.Applicant
	err = database.DB.Where("profile_code = ?", code).First(&applicant).Error
	if err != nil {
		if gorm.IsRecordNotFoundError(err) {
			err = response.WriteError(http.StatusNotFound, errors.New("applicant not found"))
		} else {
			err = response.WriteError(http.StatusInternalServerError, err)
		}
		if err != nil {
			return
		}
		return
	}

	// Find the index of the employment with the specified 'id'
	indexToRemove := -1
	for i := range applicant.Employment {
		if applicant.Employment[i].ID == id {
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
	applicant.Employment = append(applicant.Employment[:indexToRemove], applicant.Employment[indexToRemove+1:]...)

	err = database.DB.Save(&applicant).Error
	if err != nil {
		err = response.WriteError(http.StatusInternalServerError, err)
		return
	}

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

	// Find the applicant in the database
	var applicant models.Applicant
	err = database.DB.Where("profile_code = ?", code).Preload("Education").First(&applicant).Error
	if err != nil {
		if gorm.IsRecordNotFoundError(err) {
			err = response.WriteError(http.StatusNotFound, errors.New("applicant not found"))
		} else {
			err = response.WriteError(http.StatusInternalServerError, err)
		}
		if err != nil {
			return
		}
		return
	}

	result := struct {
		Education []models.Education `json:"education"`
	}{
		Education: applicant.Education,
	}

	err = response.WriteEntity(result)
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

	// Find the applicant in the database
	var applicant models.Applicant
	err = database.DB.Where("profile_code = ?", code).First(&applicant).Error
	if err != nil {
		if gorm.IsRecordNotFoundError(err) {
			err = response.WriteError(http.StatusNotFound, errors.New("applicant not found"))
		} else {
			err = response.WriteError(http.StatusInternalServerError, err)
		}
		if err != nil {
			return
		}
		return
	}

	// Generate a new ID for the new education
	newID := len(applicant.Education) + 1
	education.ID = newID

	applicant.Education = append(applicant.Education, education)

	err = database.DB.Save(&applicant).Error
	if err != nil {
		err = response.WriteError(http.StatusInternalServerError, err)
		return
	}

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

	// Find the applicant in the database
	var applicant models.Applicant
	err = database.DB.Where("profile_code = ?", code).Preload("Education").First(&applicant).Error
	if err != nil {
		if gorm.IsRecordNotFoundError(err) {
			err = response.WriteError(http.StatusNotFound, errors.New("applicant not found"))
		} else {
			err = response.WriteError(http.StatusInternalServerError, err)
		}
		if err != nil {
			return
		}
		return
	}

	// Find the index of the education with the specified 'id'
	indexToRemove := -1
	for i := range applicant.Education {
		if applicant.Education[i].ID == id {
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
	applicant.Education = append(applicant.Education[:indexToRemove], applicant.Education[indexToRemove+1:]...)

	err = database.DB.Save(&applicant).Error
	if err != nil {
		err = response.WriteError(http.StatusInternalServerError, err)
		if err != nil {
			return
		}
		return
	}

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

// GetSkillByCode handles GET requests to retrieve an applicant skill by profile code
func GetSkillByCode(request *restful.Request, response *restful.Response) {
	candidateCode := request.PathParameter("code")
	code, err := strconv.Atoi(candidateCode)
	if err != nil {
		err = response.WriteError(http.StatusBadRequest, err)
		return
	}

	// Find the applicant in the database
	var applicant models.Applicant
	err = database.DB.Where("profile_code = ?", code).Preload("Skill").First(&applicant).Error
	if err != nil {
		if gorm.IsRecordNotFoundError(err) {
			err = response.WriteError(http.StatusNotFound, errors.New("applicant not found"))
		} else {
			err = response.WriteError(http.StatusInternalServerError, err)
		}
		if err != nil {
			return
		}
		return
	}

	result := struct {
		Skill []models.Skill `json:"skill"`
	}{
		Skill: applicant.Skill,
	}

	err = response.WriteEntity(result)
	if err != nil {
		return
	}
}

// AddSkill handles POST requests to add an applicant skill by profile code
func AddSkill(request *restful.Request, response *restful.Response) {
	candidateCode := request.PathParameter("code")
	code, err := strconv.Atoi(candidateCode)
	if err != nil {
		err = response.WriteError(http.StatusBadRequest, err)
		return
	}

	skill := models.Skill{}
	err = request.ReadEntity(&skill)
	if err != nil {
		err = response.WriteError(http.StatusBadRequest, err)
		return
	}

	// Find the applicant in the database
	var applicant models.Applicant
	err = database.DB.Where("profile_code = ?", code).Preload("Skill").First(&applicant).Error
	if err != nil {
		if gorm.IsRecordNotFoundError(err) {
			err = response.WriteError(http.StatusNotFound, errors.New("applicant not found"))
		} else {
			err = response.WriteError(http.StatusInternalServerError, err)
		}
		if err != nil {
			return
		}
		return
	}

	// Generate a new ID for the new skill
	newID := len(applicant.Skill) + 1
	skill.ID = newID

	// Append the new skill to the applicant's skill list
	applicant.Skill = append(applicant.Skill, skill)

	err = database.DB.Save(&applicant).Error
	if err != nil {
		err = response.WriteError(http.StatusInternalServerError, err)
		return
	}

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

// DeleteSkill handles DELETE requests to remove an applicant skill by profile code
func DeleteSkill(request *restful.Request, response *restful.Response) {
	candidateCode := request.PathParameter("code")
	code, err := strconv.Atoi(candidateCode)
	if err != nil {
		err = response.WriteError(http.StatusBadRequest, err)
		return
	}

	// Parse the 'id' query parameter
	skillID := request.QueryParameter("id")
	id, err := strconv.Atoi(skillID)
	if err != nil {
		err = response.WriteError(http.StatusBadRequest, err)
		return
	}

	// Find the applicant in the database
	var applicant models.Applicant
	err = database.DB.Where("profile_code = ?", code).Preload("Skill").First(&applicant).Error
	if err != nil {
		if gorm.IsRecordNotFoundError(err) {
			err = response.WriteError(http.StatusNotFound, errors.New("applicant not found"))
		} else {
			err = response.WriteError(http.StatusInternalServerError, err)
		}
		if err != nil {
			return
		}
		return
	}

	// Find the index of the skill with the specified 'id'
	indexToRemove := -1
	for i := range applicant.Skill {
		if applicant.Skill[i].ID == id {
			indexToRemove = i
			break
		}
	}

	if indexToRemove == -1 {
		err = response.WriteError(http.StatusNotFound, errors.New("skill not found"))
		if err != nil {
			return
		}
		return
	}

	// Remove the skill from the applicant's skill slice
	applicant.Skill = append(applicant.Skill[:indexToRemove], applicant.Skill[indexToRemove+1:]...)

	err = database.DB.Save(&applicant).Error
	if err != nil {
		err = response.WriteError(http.StatusInternalServerError, err)
		return
	}

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
