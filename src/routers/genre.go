package routers

import (
	"history_anime/src/controllers"
	"history_anime/src/middlewares"

	"github.com/julienschmidt/httprouter"
)

func GenreRoute(genre *httprouter.Router) {

	genre.GET("/api/genre", middlewares.Logging(middlewares.OnlyLogin(controllers.GenreGetAll)))
	genre.POST("/api/genre", middlewares.Logging(middlewares.OnlyLogin(controllers.GenreAdd)))
	genre.DELETE("/api/genre/:id", middlewares.Logging(middlewares.OnlyLogin(controllers.GenreDelete)))
}
