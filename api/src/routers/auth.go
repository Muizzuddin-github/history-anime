package routers

import (
	"history_anime/api/src/controllers"

	"github.com/julienschmidt/httprouter"
)

func AuthRoute(auth *httprouter.Router) {

	auth.POST("/api/register", controllers.Register)
	auth.POST("/api/login", controllers.Login)
	auth.POST("/api/logout", controllers.Logout)
}
