package routers

import (
	"history_anime/src/controllers"
	"history_anime/src/middlewares"

	"github.com/julienschmidt/httprouter"
)

func AnimeRoute(anime *httprouter.Router) {

	anime.GET("/api/anime", middlewares.OnlyLogin(controllers.AnimeGetAll))
	anime.POST("/api/anime", middlewares.OnlyLogin(controllers.AnimeAdd))
	anime.PUT("/api/anime/:id", middlewares.OnlyLogin(controllers.AnimeUpdate))
	anime.DELETE("/api/anime/:id", middlewares.OnlyLogin(controllers.AnimeDel))

}
