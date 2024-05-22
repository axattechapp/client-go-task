package router

import (
	"client_task/pkg/controllers"
	middlewares_test "client_task/pkg/middlewares"

	"github.com/gin-gonic/gin"
)

type UserRoutes struct {
	userController controllers.UsersController
}

func NewRouteUser(userController controllers.UsersController) UserRoutes {
	return UserRoutes{userController}
}

func (cr *UserRoutes) UserRoutes(rg *gin.RouterGroup) {

	router := rg.Group("users").Use(middlewares_test.JwtAuthMiddleware())
	router.POST("/", cr.userController.CreateUser)
	router.GET("/", cr.userController.GetAllUsers)
	router.PUT("/:UserId", cr.userController.UpdateUser)
	router.GET("/:UserId", cr.userController.GetUserById)
	router.DELETE("/:UserId", cr.userController.DeleteUserById)
}

type ProfileRoutes struct {
	profileController controllers.ProfilesController
}

func NewProfileRoutes(profileController controllers.ProfilesController) ProfileRoutes {
	return ProfileRoutes{profileController}
}

func (pr *ProfileRoutes) ProfileRoutes(rg *gin.RouterGroup) {
	router := rg.Group("profiles").Use(middlewares_test.JwtAuthMiddleware()) // Assuming same authentication middleware

	router.POST("/", pr.profileController.CreateProfile)
	router.GET("/", pr.profileController.GetAllProfiles) // Implement if needed
	router.PUT("/:ProfileID", pr.profileController.UpdateProfile)
	router.GET("/:ProfileID", pr.profileController.GetProfileByID)
	router.DELETE("/:ProfileID", pr.profileController.DeleteProfileByID) // Implement if needed
}

type CareerRoutes struct {
	careerController controllers.CareersController
}

func NewCareerRoutes(careerController controllers.CareersController) CareerRoutes {
	return CareerRoutes{careerController}
}

func (pr *CareerRoutes) CareerRoutes(rg *gin.RouterGroup) {
	router := rg.Group("careers").Use(middlewares_test.JwtAuthMiddleware()) // Assuming same authentication middleware

	router.POST("/", pr.careerController.CreateCareer)
	router.GET("/all", pr.careerController.GetAllCareers)
	router.PUT("/:CareerID", pr.careerController.UpdateCareer)
	router.GET("/:CareerID", pr.careerController.GetCareerByID)
	router.GET("/user/:UserID", pr.careerController.GetCareersByUserID)
	router.DELETE("/:CareerID", pr.careerController.DeleteCareerByID)
}

type SkillsRoutes struct {
	skillsController controllers.SkillsController
}

func NewSkillsRoutes(skillsController controllers.SkillsController) SkillsRoutes {
	return SkillsRoutes{skillsController}
}
func (pr *SkillsRoutes) SkillsRoutes(rg *gin.RouterGroup) {
	router := rg.Group("skills").Use(middlewares_test.JwtAuthMiddleware()) // Assuming same authentication middleware

	router.GET("/", pr.skillsController.GetAllSkills)
	router.POST("/", pr.skillsController.CreateSkill)
	router.GET("/:UserID", pr.skillsController.GetUserSkills)

}
