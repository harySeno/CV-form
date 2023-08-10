package main

import (
	"CV-form/pkg/database"
	"CV-form/pkg/models"
	"CV-form/pkg/restapi"
	"github.com/emicklei/go-restful/v3"
	"github.com/jinzhu/gorm"
	"log"
	"net/http"
)

func main() {
	// create new web service
	webService := new(restful.WebService)
	webService.Path("/api").
		Consumes(restful.MIME_JSON).
		Produces(restful.MIME_JSON)

	db, err := database.InitDB()
	if err != nil {
		panic("Failed to connect to the database")
	}
	defer func(db *gorm.DB) {
		err := db.Close()
		if err != nil {

		}
	}(db)

	db.AutoMigrate(
		&models.Applicant{},
		&models.Employment{},
		&models.Education{},
		&models.Skill{},
	)

	// create a new route
	webService.
		Route(webService.GET("/profile").
			To(restapi.GetAllData))

	webService.
		Route(webService.GET("/profile/{code}").
			To(restapi.GetByCode))

	webService.
		Route(webService.POST("/profile").
			To(restapi.AddProfile))

	webService.
		Route(webService.PUT("/profile/{code}").
			To(restapi.UpdateProfile))

	webService.
		Route(webService.PUT("/photo/{code}").
			To(restapi.UploadPhoto))

	webService.
		Route(webService.GET("/photo/{code}").
			To(restapi.DownloadPhoto))

	webService.
		Route(webService.DELETE("/photo/{code}").
			To(restapi.DeletePhoto))

	webService.
		Route(webService.GET("/working-experience/{code}").
			To(restapi.GetExpByCode))

	webService.
		Route(webService.PUT("/working-experience/{code}").
			To(restapi.UpdateExperience))

	webService.
		Route(webService.GET("/employment/{code}").
			To(restapi.GetEmploymentByCode))

	webService.
		Route(webService.POST("/employment/{code}").
			To(restapi.AddEmployment))

	webService.
		Route(webService.DELETE("/employment/{code}").
			To(restapi.DeleteEmployment))

	webService.
		Route(webService.GET("/education/{code}").
			To(restapi.GetEducationByCode))

	webService.
		Route(webService.POST("/education/{code}").
			To(restapi.AddEducation))

	webService.
		Route(webService.DELETE("/education/{code}").
			To(restapi.DeleteEducation))

	webService.
		Route(webService.GET("/skill/{code}").
			To(restapi.GetSkillByCode))

	webService.
		Route(webService.POST("/skill/{code}").
			To(restapi.AddSkill))

	webService.
		Route(webService.DELETE("/skill/{code}").
			To(restapi.DeleteSkill))

	// create a new container
	container := restful.NewContainer()
	container.Add(webService)

	// initiate server
	log.Println(" Initiating server on http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", container))
}
