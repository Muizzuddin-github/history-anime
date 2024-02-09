package routers

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func Router() *httprouter.Router {

	routers := httprouter.New()
	routers.GET("/api", func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {

		w.Write([]byte("selamat datang"))
	})
	AuthRoute(routers)
	AnimeRoute(routers)
	GenreRoute(routers)

	return routers
}
