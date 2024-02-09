package routers

import (
	"history_anime/src/controllers"

	"github.com/julienschmidt/httprouter"
)

func AuthRoute(auth *httprouter.Router) {

	auth.POST("/api/register", controllers.Register)
	auth.POST("/api/login", controllers.Login)
	auth.POST("/api/logout", controllers.Logout)
	auth.POST("/api/forgot-password", controllers.ForgotPassword)
	auth.POST("/api/reset-password", controllers.ResetPassword)
}
