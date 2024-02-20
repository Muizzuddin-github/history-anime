package routers

import (
	"history_anime/src/controllers"
	"history_anime/src/middlewares"

	"github.com/julienschmidt/httprouter"
)

func AuthRoute(auth *httprouter.Router) {

	auth.POST("/api/register", middlewares.Logging(controllers.Register))
	auth.POST("/api/login", middlewares.Logging(controllers.Login))
	auth.POST("/api/logout", middlewares.Logging(controllers.Logout))
	auth.POST("/api/forgot-password", middlewares.Logging(controllers.ForgotPassword))
	auth.POST("/api/reset-password", middlewares.Logging(controllers.ResetPassword))
	auth.GET("/api/islogin", middlewares.Logging(controllers.IsLogin))
}
