package routers

import (
	"history_anime/src/controllers"
	"history_anime/src/middlewares"

	"github.com/julienschmidt/httprouter"
)

func GenreRoute(genre *httprouter.Router) {

	genre.GET("/api/genre", middlewares.OnlyLogin(controllers.GenreGetAll))
	genre.POST("/api/genre", middlewares.OnlyLogin(controllers.GenreAdd))
	genre.DELETE("/api/genre/:id", middlewares.OnlyLogin(controllers.GenreDelete))
}
