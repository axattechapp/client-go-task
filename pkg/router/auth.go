package router

import (
	"client_task/pkg/controllers"

	"github.com/gin-gonic/gin"
)

type UserAuthRoutes struct {
	userController controllers.UsersController
}

func NewRouteUserAuth(userController controllers.UsersController) UserAuthRoutes {
	return UserAuthRoutes{userController}
}

func (cr *UserAuthRoutes) UserAuthRoutes(rg *gin.RouterGroup) {

	router := rg.Group("auth")
	router.POST("/login", cr.userController.LoginHandler)
	router.POST("/registration", cr.userController.RegisterHandler)

}
